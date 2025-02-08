package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type tab int

const (
	statsTab tab = iota
	infoTab
	errorTab
	poolsTab
	sessionsTab
)

type model struct {
	table     table.Model
	viewport  viewport.Model
	activeTab tab
	tabs      []string
	info      string
	errors    string
	pools     string
	sessions  string
	err       error
	lastFetch time.Time
	width     int
	height    int
	config    Config
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return fetchStats(m.config) },
		func() tea.Msg { return fetchInfo(m.config) },
		func() tea.Msg { return fetchErrors(m.config) },
		func() tea.Msg { return fetchPools(m.config) },
		func() tea.Msg { return fetchSessions(m.config) },
	)
}

type (
	infoMsg    string
	errorMsg   string
	poolsMsg   string
	sessionMsg string
)

func fetchStats(cfg Config) tea.Msg {
	conn, err := net.Dial("unix", cfg.socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "show stat\n")

	var rows []table.Row
	scanner := bufio.NewScanner(conn)
	first := true

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if first {
			first = false
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 80 {
			continue
		}

		row := table.Row{
			// truncate(fields[0], 18), // Name
			// truncate(fields[1], 13), // Server
			fields[0],              // Name
			fields[1],              // Server
			fields[17],             // Status
			fields[4],              // Current Sessions
			fields[5],              // Max Sessions
			fields[7],              // Total Sessions
			formatBytes(fields[8]), // Bytes In
			formatBytes(fields[9]), // Bytes Out
			fields[13],             // Errors
			fields[18],             // Weight
		}

		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return rows
}

func parseInfoToRows(info string) []table.Row {
	var rows []table.Row
	lines := strings.Split(info, "\n")

	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split on first ':'
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		rows = append(rows, table.Row(parts))
	}

	return rows
}

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

func initializeInfoTable() table.Model {
	info_columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Value", Width: 25},
		{Title: "Description", Width: 120},
	}

	info_table := table.New(
		table.WithColumns(info_columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	info_table.SetStyles(statsTableStyles())

	return info_table
}

func initializeStatsTable() table.Model {
	stats_columns := []table.Column{
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

	stats_table := table.New(
		table.WithColumns(stats_columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	stats_table.SetStyles(statsTableStyles())

	return stats_table
}

func main() {
	cfg := Config{
		socketPath: "/var/run/haproxy/admin.sock", // default path
	}
	if len(os.Args) > 1 {
		cfg.socketPath = os.Args[1]
	}

	// Initial state
	m := model{
		table:     initializeStatsTable(),
		viewport:  viewport.New(80, 20),
		tabs:      []string{"Stats", "Info", "Errors", "Memory", "Sessions"},
		activeTab: statsTab,
		config:    cfg,
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
