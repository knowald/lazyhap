package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/knowald/lazyhap/src/views/certs"
	"github.com/knowald/lazyhap/src/views/error"
	"github.com/knowald/lazyhap/src/views/help"
	"github.com/knowald/lazyhap/src/views/info"
	"github.com/knowald/lazyhap/src/views/pools"
	"github.com/knowald/lazyhap/src/views/sessions"
	"github.com/knowald/lazyhap/src/views/stats"
	"github.com/knowald/lazyhap/src/views/threads"
)

func (m model) View() tea.View {
	var content string

	if m.showHelp {
		content = help.RenderHelp()
	} else if m.err != nil {
		content = fmt.Sprintf("\nError: %v\n\nPress q to quit\n", m.err)
	} else {
		var sb strings.Builder

		renderTabBar(&sb, m)
		sb.WriteString("\n\n")

		switch m.activeTab {
		case statsTab:
			stats.RenderTab(&sb, m, baseStyle)
		case infoTab:
			info.RenderTab(&sb, m, baseStyle)
		case errorTab:
			error.RenderTab(&sb, m, baseStyle)
		case poolsTab:
			pools.RenderTab(&sb, m, baseStyle)
		case sessionsTab:
			sessions.RenderTab(&sb, m, baseStyle)
		case certsTab:
			certs.RenderTab(&sb, m, baseStyle)
		case threadsTab:
			threads.RenderTab(&sb, m, baseStyle)
		}

		content = sb.String()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// Navigation header
func renderTabBar(sb *strings.Builder, m model) {
	renderedTabs := make([]string, len(m.tabs))
	for i, t := range m.tabs {
		if i == int(m.activeTab) {
			renderedTabs[i] = activeTabStyle.Render(t)
		} else {
			renderedTabs[i] = tabStyle.Render(t)
		}
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...))
	sb.WriteString(renderLastUpdatedTime(m))
}

// Timestamp, last updated
func renderLastUpdatedTime(m model) string {
	timestamp := timeStyle.Render(fmt.Sprintf("Updated: %s", m.lastFetch.Format("15:04:05")))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)
	hint := hintStyle.Render("Press ? for help")
	return timestamp + hint
}
