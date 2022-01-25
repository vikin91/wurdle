package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	NewGame        key.Binding
	ChangeLanguage key.Binding
	RandomWord     key.Binding
	ShowAnswer     key.Binding

	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	ChangeLanguage: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Language"),
	),
	NewGame: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "New game"),
	),
	RandomWord: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "Random word"),
	),
	ShowAnswer: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "Show answer"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Answer"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NewGame, k.RandomWord, k.ShowAnswer},     // first column
		{k.ChangeLanguage, k.Enter, k.Help, k.Quit}, // second column
	}
}
