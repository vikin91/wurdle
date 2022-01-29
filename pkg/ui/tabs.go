package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vikin91/wurdle/pkg/model"
)

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
func (ui *UI) RenderTabs(m *model.Model) string {
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

	tabsRow := renderTabs(m.ActiveTabID, m.Tabs...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, ui.Width()-lipgloss.Width(tabsRow)-2)))
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
