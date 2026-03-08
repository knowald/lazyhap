package info

import (
	"strings"

	"charm.land/lipgloss/v2"
)

type Model interface {
	TableView() string
	GetMessage() string
	FilterMode() bool
	FilterInput() string
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	sb.WriteString(baseStyle.Render(m.TableView()))
	sb.WriteString("\n")

	if m.FilterMode() {
		filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		sb.WriteString(filterStyle.Render("Filter: " + m.FilterInput() + "█"))
		sb.WriteString(" ")
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("(enter: apply  esc: clear)"))
	} else if message := m.GetMessage(); message != "" {
		sb.WriteString(message)
	} else {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("y: copy value  /: filter  ?: help  1-7: jump tabs"))
	}
}
