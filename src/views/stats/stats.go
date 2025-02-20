package stats

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Model interface {
	TableView() string
	LastFetchTime() string
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	sb.WriteString(baseStyle.Render(m.TableView()))
	sb.WriteString("\n")
	sb.WriteString("Commands: (d)isable server, (e)nable server, set (w)eight to 100")
}
