package stats

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Model interface {
	TableView() string
	LastFetchTime() string
	FilterMode() bool
	FilterInput() string
	GetTable() table.Model
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	// Get the table and render with colorization
	tbl := m.GetTable()
	sb.WriteString(baseStyle.Render(renderColorizedTable(tbl)))
	sb.WriteString("\n")

	if m.FilterMode() {
		filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		sb.WriteString(filterStyle.Render("🔍 Filter: " + m.FilterInput() + "█"))
		sb.WriteString(" ")
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("(enter: apply  esc: clear)"))
	} else {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("d: disable  e: enable  w: weight  /: filter  ?: help  1-7: jump tabs"))
	}
}

// renderColorizedTable renders the table with status and error colorization
func renderColorizedTable(tbl table.Model) string {
	// Get the base table view
	tableView := tbl.View()

	// Apply colorization to status and error columns in each row
	lines := strings.Split(tableView, "\n")
	var result strings.Builder

	for i, line := range lines {
		if i < 2 { // Header rows
			result.WriteString(line)
		} else {
			result.WriteString(colorizeTableRow(line))
		}
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// colorizeTableRow colorizes status (column 3) and errors (column 9) in a table row
func colorizeTableRow(row string) string {
	// Parse the row to find column positions while preserving spacing
	// We need to identify the status and error columns and colorize them in place

	// Find all non-space sequences (columns) and their positions
	type columnPos struct {
		start int
		end   int
		value string
	}

	var columns []columnPos
	inColumn := false
	start := 0

	for i, ch := range row {
		if ch != ' ' {
			if !inColumn {
				inColumn = true
				start = i
			}
		} else {
			if inColumn {
				columns = append(columns, columnPos{start, i, row[start:i]})
				inColumn = false
			}
		}
	}
	if inColumn {
		columns = append(columns, columnPos{start, len(row), row[start:]})
	}

	if len(columns) < 10 {
		return row
	}

	// Colorize status (column index 2) and errors (column index 8)
	// Build the result by replacing values while preserving spacing
	result := []byte(row)

	// Colorize status column (index 2)
	statusCol := columns[2]
	colorizedStatus := colorizeStatus(statusCol.value)
	result = replaceInPlace(result, statusCol.start, statusCol.end, colorizedStatus)

	// Colorize errors column (index 8) - adjust for any length change from status colorization
	errorsCol := columns[8]
	lenDiff := len(colorizedStatus) - len(statusCol.value)
	colorizedErrors := colorizeErrors(errorsCol.value)
	result = replaceInPlace(result, errorsCol.start+lenDiff, errorsCol.end+lenDiff, colorizedErrors)

	return string(result)
}

// replaceInPlace replaces a substring with a colored version, handling ANSI codes
func replaceInPlace(row []byte, start, end int, replacement string) []byte {
	before := row[:start]
	after := row[end:]

	// Build new row
	var result []byte
	result = append(result, before...)
	result = append(result, []byte(replacement)...)
	result = append(result, after...)

	return result
}

func colorizeStatus(status string) string {
	var style lipgloss.Style

	switch status {
	case "UP":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	case "DOWN":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
	case "MAINT":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow
	case "DRAIN":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // Cyan
	case "NOLB":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // Magenta
	default:
		return status
	}

	return style.Render(status)
}

func colorizeErrors(errors string) string {
	if errors == "" || errors == "0" {
		return errors
	}

	// Parse error count
	var errorCount int64
	_, err := fmt.Sscanf(errors, "%d", &errorCount)
	if err != nil {
		return errors
	}

	var style lipgloss.Style
	if errorCount == 0 {
		return errors
	} else if errorCount < 10 {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true) // Bold red
	}

	return style.Render(errors)
}
