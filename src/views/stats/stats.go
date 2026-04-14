package stats

import (
	"strings"
	"unicode/utf8"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

type Model interface {
	TableView() string
	LastFetchTime() string
	FilterMode() bool
	FilterInput() string
	GetTable() table.Model
	SortColumn() int
	SortAscending() bool
	ConfirmMode() bool
	ConfirmPrompt() string
	WeightMode() bool
	WeightInput() string
	WeightServer() string
}

func RenderTab(sb *strings.Builder, m Model, baseStyle lipgloss.Style) {
	tbl := m.GetTable()
	sb.WriteString(baseStyle.Render(renderColorizedTable(tbl)))
	sb.WriteString("\n")

	if m.ConfirmMode() {
		confirmStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
		sb.WriteString(confirmStyle.Render(m.ConfirmPrompt()))
	} else if m.WeightMode() {
		weightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		sb.WriteString(weightStyle.Render("Weight for " + m.WeightServer() + ": " + m.WeightInput() + "█"))
		sb.WriteString(" ")
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("(0-256  enter: apply  esc: cancel)"))
	} else if m.FilterMode() {
		filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		sb.WriteString(filterStyle.Render("Filter: " + m.FilterInput() + "█"))
		sb.WriteString(" ")
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		sb.WriteString(hintStyle.Render("(enter: apply  esc: clear)"))
	} else {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		hint := "d: disable  D: drain  e: enable  R: ready  w: weight  s: sort  /: filter  ?: help"
		if col := m.SortColumn(); col >= 0 {
			cols := tbl.Columns()
			if col < len(cols) {
				arrow := "▲"
				if !m.SortAscending() {
					arrow = "▼"
				}
				sortStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
				hint += "  " + sortStyle.Render(cols[col].Title+" "+arrow)
			}
		}
		sb.WriteString(hintStyle.Render(hint))
	}
}

const (
	typeColIndex   = 0
	statusColIndex = 3
	curColIndex    = 4
	maxColIndex    = 5
	limitColIndex  = 6
	errorsColIndex = 11
	cellPadding    = 2 // Padding(0, 1) adds 1 char on each side
)

// renderColorizedTable post-processes the rendered table output, injecting
// foreground-only ANSI codes into the Status and Errors columns. Uses raw SGR
// sequences (\e[XXm ... \e[39m) instead of lipgloss to avoid full-reset codes
// that would break the selected row's background highlight.
func renderColorizedTable(tbl table.Model) string {
	tableView := tbl.View()
	cols := tbl.Columns()
	if len(cols) <= errorsColIndex {
		return tableView
	}

	// Calculate the visible character offset where each column starts.
	// Each cell is rendered at col.Width + cellPadding visible chars.
	colStart := func(idx int) int {
		pos := 0
		for i := 0; i < idx; i++ {
			pos += cols[i].Width + cellPadding
		}
		return pos
	}
	colEnd := func(idx int) int {
		return colStart(idx) + cols[idx].Width + cellPadding
	}

	typeStart, typeEnd := colStart(typeColIndex), colEnd(typeColIndex)
	statusStart, statusEnd := colStart(statusColIndex), colEnd(statusColIndex)
	curStart, curEnd := colStart(curColIndex), colEnd(curColIndex)
	maxStart, maxEnd := colStart(maxColIndex), colEnd(maxColIndex)
	limitStart, limitEnd := colStart(limitColIndex), colEnd(limitColIndex)
	errorsStart, errorsEnd := colStart(errorsColIndex), colEnd(errorsColIndex)

	lines := strings.Split(tableView, "\n")

	// headersView() produces 2 lines (header text + border from BorderBottom),
	// then View() adds "\n" before viewport content. Data rows start at index 2.
	const dataStart = 2

	var result strings.Builder
	for i, line := range lines {
		if i >= dataStart {
			// Extract the limit value for this row to color Cur and Max relative to it
			limitVal := extractCellValue(line, limitStart, limitEnd)
			curColorFn := sessionsColor(limitVal)
			maxColorFn := sessionsColor(limitVal)

			line = colorizeCellRange(line, statusStart, statusEnd, statusColor)
			line = colorizeCellRange(line, curStart, curEnd, curColorFn)
			line = colorizeCellRange(line, maxStart, maxEnd, maxColorFn)
			line = colorizeCellRange(line, limitStart, limitEnd, limitColor(limitVal))
			line = colorizeCellRange(line, errorsStart, errorsEnd, errorsColor)
			line = boldAggregateRow(line, typeStart, typeEnd)
		}
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// colorizeCellRange injects a foreground ANSI color into the cell at the given
// visible-character range. colorFn inspects the trimmed cell value and returns
// an ANSI SGR color code (e.g. "32" for green), or "" to skip colorization.
func colorizeCellRange(line string, visStart, visEnd int, colorFn func(string) string) string {
	// Walk the line tracking visible character count vs byte position,
	// skipping over ANSI escape sequences.
	byteStart := visibleToBytePos(line, visStart)
	byteEnd := visibleToBytePos(line, visEnd)

	if byteStart >= len(line) {
		return line
	}
	if byteEnd > len(line) {
		byteEnd = len(line)
	}

	cell := line[byteStart:byteEnd]
	plain := stripANSI(cell)
	value := strings.TrimSpace(plain)
	if value == "" {
		return line
	}

	sgr := colorFn(value)
	if sgr == "" {
		return line
	}

	// Find the value within the plain cell text and convert the byte offset
	// to a rune offset, since visibleToBytePos counts runes not bytes.
	valByteOffset := strings.Index(plain, value)
	if valByteOffset < 0 {
		return line
	}
	valRuneOffset := utf8.RuneCountInString(plain[:valByteOffset])
	valRuneLen := utf8.RuneCountInString(value)

	// Map the value's rune start/end to byte positions within the cell
	valByteStart := byteStart + visibleToBytePos(cell, valRuneOffset)
	valByteEnd := byteStart + visibleToBytePos(cell, valRuneOffset+valRuneLen)

	// Inject foreground-only SGR: \e[XXm before, \e[39m after (default fg reset).
	// This avoids \e[0m which would kill the Selected row's background.
	var out strings.Builder
	out.WriteString(line[:valByteStart])
	out.WriteString(sgr)
	out.WriteString(line[valByteStart:valByteEnd])
	out.WriteString("\x1b[39m")
	out.WriteString(line[valByteEnd:])
	return out.String()
}

// statusColor returns an ANSI SGR foreground sequence for the given status, or "".
func statusColor(status string) string {
	switch status {
	case "UP":
		return "\x1b[32m" // green
	case "DOWN":
		return "\x1b[31m" // red
	case "MAINT":
		return "\x1b[33m" // yellow
	case "DRAIN":
		return "\x1b[36m" // cyan
	case "NOLB":
		return "\x1b[35m" // magenta
	default:
		return ""
	}
}

// errorsColor returns an ANSI SGR foreground sequence for the given error count, or "".
func errorsColor(errors string) string {
	if errors == "" || errors == "0" {
		return ""
	}
	var n int64
	for _, c := range errors {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		} else {
			return ""
		}
	}
	if n == 0 {
		return ""
	}
	if n < 10 {
		return "\x1b[33m" // yellow
	}
	return "\x1b[1;31m" // bold red
}

// boldAggregateRow makes FE/BE rows bold by checking the Type column value.
func boldAggregateRow(line string, typeStart, typeEnd int) string {
	byteStart := visibleToBytePos(line, typeStart)
	byteEnd := visibleToBytePos(line, typeEnd)
	if byteStart >= len(line) {
		return line
	}
	if byteEnd > len(line) {
		byteEnd = len(line)
	}
	cell := line[byteStart:byteEnd]
	plain := strings.TrimSpace(stripANSI(cell))
	if plain == "→FE" || plain == "∑BE" {
		return "\x1b[1m" + line + "\x1b[22m"
	}
	return line
}

// visibleToBytePos maps a visible rune position (ignoring ANSI escapes)
// to the corresponding byte position in s.
func visibleToBytePos(s string, visPos int) int {
	visible := 0
	inEscape := false
	for i := 0; i < len(s); {
		if visible == visPos && !inEscape {
			return i
		}
		if s[i] == '\x1b' {
			inEscape = true
			i++
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			i++
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		visible++
	}
	return len(s)
}

// stripANSI removes ANSI escape codes from a string.
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

// extractCellValue reads the plain text value from a cell at the given visible range.
func extractCellValue(line string, visStart, visEnd int) string {
	byteStart := visibleToBytePos(line, visStart)
	byteEnd := visibleToBytePos(line, visEnd)
	if byteStart >= len(line) {
		return ""
	}
	if byteEnd > len(line) {
		byteEnd = len(line)
	}
	return strings.TrimSpace(stripANSI(line[byteStart:byteEnd]))
}

// parseNum parses a numeric string, returning 0 if empty or invalid.
func parseNum(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		} else {
			return 0
		}
	}
	return n
}

// sessionsColor returns a color function that compares a session count against
// the limit (maxconn). Returns red at 100%, yellow at >=80%, no color otherwise.
func sessionsColor(limit string) func(string) string {
	limitN := parseNum(limit)
	return func(value string) string {
		if limitN <= 0 {
			return ""
		}
		n := parseNum(value)
		if n <= 0 {
			return ""
		}
		if n >= limitN {
			return "\x1b[1;31m" // bold red - at or over limit
		}
		if n*100/limitN >= 80 {
			return "\x1b[33m" // yellow - approaching limit
		}
		return ""
	}
}

// limitColor returns a color function for the Limit column itself.
// Shows the limit in dim when set, invisible otherwise.
func limitColor(limit string) func(string) string {
	return func(value string) string {
		if limit == "" || limit == "0" {
			return ""
		}
		return "\x1b[2m" // dim
	}
}
