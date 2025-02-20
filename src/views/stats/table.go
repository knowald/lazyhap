package stats

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func InitializeTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Server", Width: 25},
		{Title: "Status", Width: 8},
		{Title: "Cur Sess", Width: 10},
		{Title: "Max Sess", Width: 10},
		{Title: "Tot Sess", Width: 10},
		{Title: "Bytes In", Width: 12},
		{Title: "Bytes Out", Width: 12},
		{Title: "Errors", Width: 8},
		{Title: "Weight", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	t.SetStyles(tableStyles())
	return t
}

func tableStyles() table.Styles {
	stats_table_styles := table.DefaultStyles()
	stats_table_styles.Header = stats_table_styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("205"))
	stats_table_styles.Selected = stats_table_styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	return stats_table_styles
}
