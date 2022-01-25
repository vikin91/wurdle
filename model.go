package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	dict       map[string]*Dictionary
	gameID     int
	gamesWon   uint
	secret     string
	won        bool
	textInput  textinput.Model
	statusText string

	answers    []string
	answerEval []string
	err        error

	tabs        []string
	activeTabID int
	lang        string

	// help
	keys keyMap
	help help.Model
}

func (m *model) CurrentDict() *Dictionary {
	if d, ok := m.dict[m.lang]; ok {
		return d
	}
	m.lang = "en"
	return m.dict["en"]
}

func newModel() *model {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	ti := textinput.New()
	ti.Placeholder = "word"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	dictMap := make(map[string]*Dictionary)
	dictMap["en"] = NewDictionary("en")

	m := model{
		dict:        dictMap,
		lang:        "en",
		gameID:      0,
		gamesWon:    0,
		won:         false,
		secret:      "",
		statusText:  "Ready to play",
		textInput:   ti,
		tabs:        []string{"EN", "DE", "PL"},
		activeTabID: 0,
		answers:     make([]string, 0),
		answerEval:  make([]string, 0),
		help:        help.New(),
		keys:        keys,
	}
	m.NewGame()
	return &m
}
