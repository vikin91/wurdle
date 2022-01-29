package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/vikin91/wurdle/pkg/action"
	"github.com/vikin91/wurdle/pkg/game"
	"github.com/vikin91/wurdle/pkg/model"
)

func pipe(s string) string {
	if len(s) > 0 {
		return s + " | "
	}
	return " "
}

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) Render(m *model.Model, help *help.Model) string {
	var b strings.Builder
	var answerStyle = lipgloss.NewStyle().Width(ui.Width()).Align(lipgloss.Center)

	b.WriteString(ui.RenderTabs(m))
	b.WriteString("\n")

	activeGame, err := m.ActiveGame()
	if err != nil {
		m.LogError(err)
	}

	numAnswers := 0
	if activeGame != nil {
		numAnswers = len(activeGame.Answers())
		for i, answer := range activeGame.Answers() {
			if i > m.MaxNumAttempts+1 {
				break
			}
			answerRender := ui.RenderAnswer(activeGame.RevealSecret(), answer)
			b.WriteString(answerStyle.Render(answerRender))
			b.WriteString("\n")
		}
	}
	for i := numAnswers; i < m.MaxNumAttempts; i++ {
		b.WriteString(answerStyle.Render(renderPlaceholder("______")))
		b.WriteString("\n")
	}
	b.WriteString("\n\n")
	b.WriteString("Word: ")
	b.WriteString(m.TextInput.View())
	b.WriteString("\n\n")

	b.WriteString(ui.renderStatusBar(m, activeGame, help))
	return b.String()
}

func (ui *UI) Width() int {
	return 96
}

func (ui *UI) renderStatusBar(m *model.Model, game *game.Game, help *help.Model) string {
	var statusBarStyle = lipgloss.NewStyle().
		Width(ui.Width()).
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
	var statusErrorStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FF0000")).
		MarginLeft(2)
	var statusLogStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#8888FF")).
		MarginLeft(2)

	numErrorsToDisplay := min(len(m.Errors), m.MaxNumErrors)
	errMsgs := make([]string, numErrorsToDisplay)
	for i := 0; i < numErrorsToDisplay; i++ {
		errMsg := fmt.Sprintf("Error[%d]: %s", i, m.Errors[i].Error())
		errMsgs[i] = statusErrorStyle.Render(errMsg)
	}

	numLogsToDisplay := min(len(m.Logs), m.MaxNumLines)
	logMsgs := make([]string, numLogsToDisplay)
	for i := 0; i < numLogsToDisplay; i++ {
		errMsg := fmt.Sprintf("Log[%d]: %s", i, m.Logs[i])
		logMsgs[i] = statusLogStyle.Render(errMsg)
	}

	gamesWonTxt := "Games won: 0 of 0"
	attemptsTxt := "Attempts: 0 of 0"
	if game != nil {
		gamesWonTxt = fmt.Sprintf("Games won: %d of %d", game.NumGamesWon, m.GameID())
		attemptsTxt = fmt.Sprintf("Attempts: %d of %d", len(game.Answers()), m.MaxNumAttempts)
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Bottom,
		pipe(gamesWonTxt),
		pipe(attemptsTxt),
		pipe(m.StatusText),
	)
	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		statusBarStyle.Width(ui.Width()).Render(bar),
		lipgloss.JoinVertical(lipgloss.Top, append(errMsgs, logMsgs...)...),
		help.View(action.Keys),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
