package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

func InitializeTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Value", Width: 25},
		{Title: "Description", Width: 120},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	t.SetStyles(tableStyles())
	return t
}

func ParseInfoToRows(info string) []table.Row {
	var rows []table.Row
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		trimmedParts := make([]string, len(parts))
		for i, part := range parts {
			trimmedParts[i] = strings.TrimSpace(part)
		}
		rows = append(rows, table.Row(trimmedParts))
	}
	return rows
}

func tableStyles() table.Styles {
	// Move your statsTableStyles() function content here
	// You'll need to add the actual styles implementation
	return table.DefaultStyles()
}
