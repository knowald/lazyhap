package threads

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/knowald/lazyhap/src/colorize"
)

type Model interface {
	ThreadsView() string
	GetViewport() viewport.Model
	ViewportFilterMode() bool
	ViewportFilterInput() string
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	viewport := m.GetViewport()
	content := colorize.ColorizeThreadOutput(m.ThreadsView())
	viewport.SetContent(content)
	sb.WriteString(baseStyle.Render(viewport.View()))
	sb.WriteString("\n")

	if m.ViewportFilterMode() {
		filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		sb.WriteString(filterStyle.Render("Filter: " + m.ViewportFilterInput() + "█"))
		sb.WriteString(" ")
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("(enter: apply  esc: clear)"))
	} else {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("j/k: scroll  /: filter  ?: help  1-9: jump tabs"))
	}
}
