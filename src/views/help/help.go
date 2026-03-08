package help

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var helpText = `LazyHAP - HAProxy TUI Tool

NAVIGATION
  tab, right, l     Next tab
  shift+tab, left, h Previous tab
  1-9               Jump to tab by number
  j, down           Move down in table/list
  k, up             Move up in table/list
  g                 Go to top
  G                 Go to bottom
  r                 Refresh current tab
  y                 Copy selected value to clipboard
  ?                 Toggle this help screen
  q, esc, ctrl+c    Quit

STATS TAB (Tab 1)
  /                 Start filtering (type to search)
  d                 Disable selected server
  D                 Drain selected server
  e                 Enable selected server
  R                 Set server state to ready
  w                 Set server weight (input popup)
  x                 Kill sessions (with confirmation)
  c                 Clear all counters
  s                 Cycle sort column (asc/desc)

INFO TAB (Tab 2)
  /                 Start filtering (type to search)

FILTER MODE (All Tabs)
  Type to search    Filter servers/backends
  Enter             Apply filter and exit mode
  Esc               Clear filter and exit mode
  Backspace         Delete last character

STATUS COLORS
  Green (UP)        Server is healthy and active
  Red (DOWN)        Server is down or unreachable
  Yellow (MAINT)    Server in maintenance mode
  Cyan (DRAIN)      Server is draining connections
  Magenta (NOLB)    Server not load-balanced

ERROR COLORS
  Yellow            Low error count (< 10)
  Bold Red          High error count (≥ 10)

TABS
  1. Stats          Server statistics and control
  2. Info           HAProxy configuration information
  3. Errors         Error logs
  4. Memory         Memory pool statistics
  5. Sessions       Active sessions
  6. Certs          SSL certificate information
  7. Threads        Thread information
  8. Activity        System activity metrics
  9. Events          Event sinks and logs

Press ? or q to close this help screen`

func RenderHelp() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(80)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center)

	lines := strings.Split(helpText, "\n")
	var formatted strings.Builder

	for i, line := range lines {
		if i == 0 {
			// Title
			formatted.WriteString(titleStyle.Render(line))
		} else if strings.HasPrefix(line, "  ") {
			// Indented lines (keyboard shortcuts)
			parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(parts) == 2 {
				keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
				formatted.WriteString("  ")
				formatted.WriteString(keyStyle.Render(parts[0]))
				formatted.WriteString(" ")
				formatted.WriteString(parts[1])
			} else {
				formatted.WriteString(line)
			}
		} else if line != "" && !strings.HasPrefix(line, " ") {
			// Section headers
			headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
			formatted.WriteString(headerStyle.Render(line))
		} else {
			formatted.WriteString(line)
		}
		if i < len(lines)-1 {
			formatted.WriteString("\n")
		}
	}

	return style.Render(formatted.String())
}
