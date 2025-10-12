package certs

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model interface {
	CertsView() string
	GetViewport() viewport.Model
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	viewport := m.GetViewport()
	viewport.SetContent(m.CertsView())
	sb.WriteString(baseStyle.Render(viewport.View()))
	sb.WriteString("\n")
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	sb.WriteString(hintStyle.Render("j/k: scroll  ?: help  1-7: jump tabs"))
}
