package info

import (
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

const defaultTableHeight = 20

func InitializeTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Value", Width: 25},
		{Title: "Description", Width: 120},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(defaultTableHeight),
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
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("205"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
}
