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
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	viewport := m.GetViewport()
	content := colorize.ColorizeThreadOutput(m.ThreadsView())
	viewport.SetContent(content)
	sb.WriteString(baseStyle.Render(viewport.View()))
	sb.WriteString("\n")
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	sb.WriteString(hintStyle.Render("j/k: scroll  ?: help  1-7: jump tabs"))
}
