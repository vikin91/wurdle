package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func pipe(s string) string {
	if len(s) > 0 {
		return s + " | "
	}
	return " "
}
func statusBar(m *model) string {
	var statusBarStyle = lipgloss.NewStyle().
		Width(m.help.Width).
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
	var statusErrorStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		MarginRight(1)

	errMsg := ""
	if m.err != nil {
		errMsg = statusErrorStyle.Render(m.err.Error())
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Bottom,
		fmt.Sprintf("Games won: %d of %d | ", m.gamesWon, m.gameID),
		fmt.Sprintf("Attempts: %d of %d | ", len(m.answers), NumAttempts),
		pipe(m.statusText),
		pipe(errMsg),
	)
	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		statusBarStyle.Width(width).Render(bar),
		m.help.View(m.keys),
	)
}

var activeTabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      " ",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┘",
	BottomRight: "└",
}

var tabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┴",
	BottomRight: "┴",
}

var highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
var tab = lipgloss.NewStyle().
	Border(tabBorder, true).
	BorderForeground(highlight).
	Padding(0, 1)

// Tabs
func displayTabs(m *model) string {
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	special := lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider := lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url := lipgloss.NewStyle().Foreground(special).Render("Wurdly\n")

	tabGap := tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	tabsRow := renderTabs(m.activeTabID, m.tabs...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(tabsRow)-2)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, url, tabsRow, gap, divider)
}

func renderTabs(activeID int, tabs ...string) string {
	activeTab := tab.Copy().Border(activeTabBorder, true)

	styledTabs := make([]string, len(tabs))
	for i, t := range tabs {
		if i == activeID {
			styledTabs = append(styledTabs, activeTab.Render(t))
		} else {
			styledTabs = append(styledTabs, tab.Render(t))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, styledTabs...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
