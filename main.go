package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

const (
	width       = 96
	NumAttempts = 6
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var (
		daemonMode bool
		showHelp   bool
		opts       []tea.ProgramOption
	)

	flag.BoolVar(&daemonMode, "d", false, "run as a daemon")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if daemonMode || !isatty.IsTerminal(os.Stdout.Fd()) {
		// If we're in daemon mode don't render the TUI
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	} else {
		// If we're in TUI mode, discard log output
		log.SetOutput(ioutil.Discard)
	}

	p := tea.NewProgram(newModel(), opts...)
	if err := p.Start(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}

func (m *model) processAnswer(word string) {
	evaluated, win, err := m.evaluateAndRenderAnswer(m.secret, strings.ToLower(word))
	m.err = err
	m.textInput.SetValue("")
	if err != nil {
		return
	}
	m.won = win
	if win {
		m.statusText = "You won! Press 1 to start new game"
	} else {
		m.statusText = ""
	}
	m.answerEval = append(m.answerEval, evaluated)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		// NEW GAME
		case key.Matches(msg, m.keys.NewGame):
			m.NewGame()

			// RANDOM
		case key.Matches(msg, m.keys.RandomWord):
			if !m.won {
				m.textInput.SetValue(m.CurrentDict().GetSecret())
			}

			// GIVE UP
		case key.Matches(msg, m.keys.ShowAnswer):
			if !m.won {
				m.textInput.SetValue(m.secret)
			}

			// LANGUAGE
		case key.Matches(msg, m.keys.ChangeLanguage):
			m.activeTabID++
			if m.activeTabID >= len(m.tabs) {
				m.activeTabID = 0
			}
			m.textInput.SetValue("")

			// HELP
		case key.Matches(msg, m.keys.Help):
			inputState := m.textInput.Value()
			m.help.ShowAll = !m.help.ShowAll
			m.textInput.SetValue(inputState)

			// QUIT
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

			// ENTER
		case key.Matches(msg, m.keys.Enter):
			if !m.won {
				m.processAnswer(m.textInput.Value())
			}

		default:
			if !m.won {
				m.textInput, cmd = m.textInput.Update(msg)
			}
		}
	}
	return m, cmd
}

func (m *model) NewGame() {
	m.gameID++
	m.textInput.SetValue("")
	m.secret = m.CurrentDict().GetSecret()
	m.answers = make([]string, 0)
	m.answerEval = make([]string, 0)
	if m.won {
		m.gamesWon++
	}
	m.won = false
	m.textInput.Placeholder = ""
}

func (m *model) View() string {
	var b strings.Builder
	var answerStyle = lipgloss.NewStyle().Width(m.help.Width).Align(lipgloss.Center)

	b.WriteString(displayTabs(m))
	b.WriteString("\n")

	numAnswers := 0
	for i, eval := range m.answerEval {
		b.WriteString(answerStyle.Render(eval))
		b.WriteString("\n")
		numAnswers = i
	}
	for i := numAnswers; i < NumAttempts+1; i++ {
		b.WriteString(answerStyle.Render(renderPlaceholder("______")))
		b.WriteString("\n")
	}
	b.WriteString("\n\n")
	b.WriteString("Word: ")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	b.WriteString(statusBar(m))

	return b.String()
}
