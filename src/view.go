package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\nPress q to quit\n", m.err)
	}

	var sb strings.Builder

	renderedTabs := []string{}
	for i, t := range m.tabs {
		if i == int(m.activeTab) {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...))
	sb.WriteString("\n\n")

	switch m.activeTab {
	case statsTab:
		sb.WriteString(renderLastUpdatedTime(m))
		sb.WriteString("\n")
		sb.WriteString(baseStyle.Render(m.table.View()))
		sb.WriteString("\n")
		sb.WriteString("Commands: (d)isable server, (e)nable server, set (w)eight to 100")

	case infoTab:
		sb.WriteString(baseStyle.Render(m.table.View()))
		sb.WriteString("\n")
		sb.WriteString("Commands: (y) to yank value to clipboard")

	case errorTab:
		m.viewport.SetContent(m.errors)
		sb.WriteString(baseStyle.Render(m.viewport.View()))

	case poolsTab:
		m.viewport.SetContent(m.pools)
		sb.WriteString(baseStyle.Render(m.viewport.View()))

	case sessionsTab:
		m.viewport.SetContent(m.sessions)
		sb.WriteString(baseStyle.Render(m.viewport.View()))
	}

	return sb.String()
}

func renderLastUpdatedTime(m model) (Rendered string) {
	return timeStyle.Render(fmt.Sprintf("Last updated: %s", m.lastFetch.Format("15:04:05")))
}
