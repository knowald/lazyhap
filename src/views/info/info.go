package info

import (
	"strings"

	"charm.land/lipgloss/v2"
)

type Model interface {
	TableView() string
	GetMessage() string
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	sb.WriteString(baseStyle.Render(m.TableView()))
	sb.WriteString("\n")
	if message := m.GetMessage(); message != "" {
		sb.WriteString(message)
	} else {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("y: copy value  ?: help  1-7: jump tabs"))
	}
}
