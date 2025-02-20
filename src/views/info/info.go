package info

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
		sb.WriteString("Commands: (y) to yank value to clipboard")
	}
}
