package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vikin91/wurdle/pkg/action"
	"github.com/vikin91/wurdle/pkg/game"
)

type Model struct {
	Games          map[string]*game.Game
	gameID         int
	MaxNumAttempts int

	ActiveTabID int
	Tabs        []string

	MaxNumErrors int
	Errors       []error

	MaxNumLines int
	Logs        []string

	TextInput    textinput.Model
	selectedLang string
	StatusText   string

	// help
	keys action.KeyMap
}

func (m *Model) LastError() error {
	if len(m.Errors) > 0 {
		return m.Errors[len(m.Errors)-1]
	}
	return nil
}

func (m *Model) ActiveGame() (*game.Game, error) {
	if game, ok := m.Games[m.selectedLang]; ok {
		return game, nil
	}
	return nil, fmt.Errorf("no active game")
}

func (m *Model) MustActiveGame() *game.Game {
	if game, ok := m.Games[m.selectedLang]; ok {
		return game
	}
	panic("no active game")
}

func (m *Model) GameID() int {
	return m.gameID
}

func NewModel(maxNumAttempts int) *Model {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	ti := textinput.New()
	ti.Placeholder = "word"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	m := Model{
		Games:          make(map[string]*game.Game),
		gameID:         0,
		MaxNumAttempts: maxNumAttempts,
		ActiveTabID:    0,
		Tabs:           make([]string, 0),
		MaxNumErrors:   5,
		Errors:         make([]error, 0),
		MaxNumLines:    5,
		Logs:           make([]string, 0),
		TextInput:      ti,
		StatusText:     "Ready to play",
		keys:           action.Keys,
	}
	return &m
}

func (m *Model) LogError(err error) {
	m.Errors = append([]error{err}, m.Errors...)
}

func (m *Model) LogInfof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	m.Logs = append([]string{s}, m.Logs...)
}

func (m *Model) AddGames(langs ...string) {
	games := make([]*game.Game, len(langs))
	for i, lan := range langs {
		g, err := game.CreateGame(lan, 6, m)
		if err != nil {
			m.LogError(err)
		}
		games[i] = g
	}
	m.addGames(games...)
}

func (m *Model) addGames(games ...*game.Game) {
	for _, g := range games {
		m.Games[g.Language()] = g
		m.Tabs = append(m.Tabs, g.Language())
		if m.selectedLang == "" {
			m.selectedLang = g.Language()
		}
	}
}

func (m *Model) processAnswer(game *game.Game) {
	word := m.TextInput.Value()
	if game.Won() {
		game.NumGamesWon++
	} else if game.CanAttempt(m.MaxNumAttempts) {
		err := game.RecordAnswer(strings.ToLower(word))
		m.TextInput.SetValue("")
		if err != nil {
			m.LogError(err)
			return
		}
	}
	if game.Lost(m.MaxNumAttempts) {
		m.StatusText = fmt.Sprintf("Lost! The word was: '%s'", game.RevealSecret())
		m.LogInfof("No more attempts. Press 1 to start new game")
	}
	if game.Won() {
		m.StatusText = "Correct! You won!"
		m.LogInfof(word + ", yeah, that was it!")
	} else {
		m.LogInfof(word + ", not bad, try next word")
		m.StatusText = ""
	}
}

func (m *Model) GetActiveLang() string {
	return m.selectedLang
}

func (m *Model) HandleCommand(msg tea.Msg, help *help.Model) tea.Cmd {
	var cmd tea.Cmd
	activeGame, err := m.ActiveGame()
	if err != nil {
		m.LogError(err)
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		// NEW GAME
		case key.Matches(msg, m.keys.NewGame) && activeGame != nil:
			activeGame.NewGame()

			// RANDOM
		case key.Matches(msg, m.keys.RandomWord) && activeGame != nil:
			if !activeGame.Won() && activeGame.CanAttempt(m.MaxNumAttempts) {
				m.TextInput.SetValue(activeGame.Dictionary.GetSecret())
			}

			// GIVE UP
		case key.Matches(msg, m.keys.ShowAnswer) && activeGame != nil:
			if !activeGame.Won() && activeGame.CanAttempt(m.MaxNumAttempts) {
				m.TextInput.SetValue(activeGame.RevealSecret())
			}

			// LANGUAGE
		case key.Matches(msg, m.keys.ChangeLanguage):
			m.ActiveTabID++
			if m.ActiveTabID >= len(m.Tabs) {
				m.ActiveTabID = 0
			}
			m.TextInput.SetValue("")
			m.selectedLang = m.Tabs[m.ActiveTabID]

			// HELP
		case key.Matches(msg, m.keys.Help):
			inputState := m.TextInput.Value()
			help.ShowAll = !help.ShowAll
			m.TextInput.SetValue(inputState)

			// QUIT
		case key.Matches(msg, m.keys.Quit):
			return tea.Quit

			// ENTER
		case key.Matches(msg, m.keys.Enter) && activeGame != nil:
			m.processAnswer(activeGame)

		default:
			if activeGame != nil && !activeGame.Won() && activeGame.CanAttempt(m.MaxNumAttempts) {
				m.TextInput, cmd = m.TextInput.Update(msg)
			}
		}
	}
	return cmd
}
