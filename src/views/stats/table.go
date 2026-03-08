package stats

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

const defaultTableHeight = 20

func InitializeTable() table.Model {
	columns := []table.Column{
		{Title: "Type", Width: 6},
		{Title: "Name", Width: 38},
		{Title: "Server", Width: 24},
		{Title: "Status", Width: 10},
		{Title: "Cur", Width: 6},
		{Title: "Max", Width: 6},
		{Title: "Total", Width: 8},
		{Title: "Bytes In", Width: 10},
		{Title: "Bytes Out", Width: 10},
		{Title: "Rate/s", Width: 7},
		{Title: "Errors", Width: 7},
		{Title: "Weight", Width: 7},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(defaultTableHeight),
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
