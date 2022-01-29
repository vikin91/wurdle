package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vikin91/wurdle/pkg/model"
	"github.com/vikin91/wurdle/pkg/ui"
)

type State struct {
	model *model.Model
	ui    *ui.UI
	help  *help.Model
}

func (s *State) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (s *State) View() string {
	return s.ui.Render(s.model, s.help)
}

func (s *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := s.model.HandleCommand(msg, s.help)
	return s, cmd
}
