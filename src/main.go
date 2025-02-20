package main

import (
	"bufio"
	"fmt"
	"lazyhap/src/views/stats"
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
	certsTab
)

type model struct {
	table     table.Model
	message   string
	viewport  viewport.Model
	activeTab tab
	tabs      []string
	info      string
	errors    string
	pools     string
	certs     string
	sessions  string
	err       error
	lastFetch time.Time
	width     int
	height    int
	config    Config
}

type (
	infoMsg    string
	errorMsg   string
	poolsMsg   string
	sessionMsg string
	certsMsg   string
)

type clearMessageMsg struct{}

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
		// skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		// split on first ':'
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		// trim the spaces
		trimmedParts := make([]string, len(parts))
		for i, part := range parts {
			trimmedParts[i] = strings.TrimSpace(part)
		}
		rows = append(rows, table.Row(trimmedParts))
	}
	return rows
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
		table:     stats.InitializeTable(),
		viewport:  viewport.New(80, 20),
		tabs:      []string{"Stats", "Info", "Errors", "Memory", "Sessions", "Certs"},
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
