package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	letterBase              = lipgloss.NewStyle().Align(lipgloss.Center).Padding(0, 1).UnsetForeground().UnsetBackground()
	letterStyleNormal       = letterBase.Copy().Bold(true).Background(lipgloss.Color("#000000")).Foreground(lipgloss.NoColor{})
	letterStyleGoodPosition = letterBase.Copy().Bold(true).Background(lipgloss.Color("#00FF00")).Foreground(lipgloss.Color("#000000"))
	letterStyleGoodLetter   = letterBase.Copy().Bold(true).Background(lipgloss.Color("#FFFF00")).Foreground(lipgloss.Color("#000000"))
)

func renderPlaceholder(s string) string {
	return renderWithPadding(s, letterBase)
}

func renderWithPadding(s string, style lipgloss.Style) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteString(style.Render(string(r)))
	}
	return b.String()
}

func (m *model) evaluateAndRenderAnswer(secret, answer string) (styled string, won bool, err error) {
	if len(answer) > 6 {
		return "", false, fmt.Errorf("word too long")
	}
	if len(answer) < 6 {
		return "", false, fmt.Errorf("word too short")
	}
	if m.CurrentDict().IsWordValid(answer) {
		m.answers = append(m.answers, answer)
	} else {
		return "", false, fmt.Errorf("'%s' is not a word", answer)
	}
	win := false
	if answer == secret {
		return renderWithPadding(answer, letterStyleGoodPosition), true, nil
	}
	alreadyYellow := make(map[rune]int)
	for _, r := range secret {
		alreadyYellow[r] = alreadyYellow[r] + 1
	}

	var b strings.Builder
	for i, r := range answer {
		if r == []rune(secret)[i] {
			b.WriteString(letterStyleGoodPosition.Render(string(r)))
		} else if strings.Count(secret, string(r)) > 0 && alreadyYellow[r] > 0 {
			alreadyYellow[r] = alreadyYellow[r] - 1
			b.WriteString(letterStyleGoodLetter.Render(string(r)))
		} else {
			b.WriteString(letterStyleNormal.Render(string(r)))
		}
	}
	return b.String(), win, nil
}
