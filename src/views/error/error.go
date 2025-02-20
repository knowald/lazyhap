package error

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model interface {
	ErrorView() string
	GetViewport() viewport.Model
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	viewport := m.GetViewport()
	viewport.SetContent(m.ErrorView())
	sb.WriteString(baseStyle.Render(viewport.View()))
}
