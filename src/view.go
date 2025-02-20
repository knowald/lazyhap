package main

import (
	"fmt"
	"lazyhap/src/views/certs"
	"lazyhap/src/views/error"
	"lazyhap/src/views/info"
	"lazyhap/src/views/stats"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\nPress q to quit\n", m.err)
	}

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
		// TODO: Separate into package
		renderPoolsTab(&sb, m)
	case sessionsTab:
		// TODO: Separate into package
		renderSessionsTab(&sb, m)
	case certsTab:
		certs.RenderTab(&sb, m, baseStyle)
	}

	return sb.String()
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

// Tabs
func renderPoolsTab(sb *strings.Builder, m model) {
	m.viewport.SetContent(m.pools)
	sb.WriteString(baseStyle.Render(m.viewport.View()))
}

func renderSessionsTab(sb *strings.Builder, m model) {
	m.viewport.SetContent(m.sessions)
	sb.WriteString(baseStyle.Render(m.viewport.View()))
}

// Timestamp, last updated
func renderLastUpdatedTime(m model) string {
	return timeStyle.Render(fmt.Sprintf("Last updated: %s", m.lastFetch.Format("15:04:05")))
}
