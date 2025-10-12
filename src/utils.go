package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	exec "os/exec"
)

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func formatBytes(bytes string) string {
	b := stringToInt(bytes)
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func stringToInt(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// colorizeStatus returns a color-coded status string
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
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Gray
	}

	return style.Render(status)
}

// colorizeErrors returns a color-coded error count
func colorizeErrors(errors string) string {
	if errors == "" || errors == "0" {
		return errors
	}

	// Parse error count
	var errorCount int64
	fmt.Sscanf(errors, "%d", &errorCount)

	var style lipgloss.Style
	if errorCount == 0 {
		style = lipgloss.NewStyle() // Default
	} else if errorCount < 10 {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow for low errors
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true) // Bold red for high errors
	}

	return style.Render(errors)
}

// filterRows filters table rows based on search query
func filterRows(rows []table.Row, query string) []table.Row {
	if query == "" {
		return rows
	}

	query = strings.ToLower(query)
	var filtered []table.Row

	for _, row := range rows {
		// Search in name, server, and status columns
		for _, cell := range row {
			// Remove ANSI color codes for searching
			plainCell := stripANSI(cell)
			if strings.Contains(strings.ToLower(plainCell), query) {
				filtered = append(filtered, row)
				break
			}
		}
	}

	return filtered
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	// Simple ANSI stripper - looks for ESC sequences
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
