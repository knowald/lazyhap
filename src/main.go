package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knowald/lazyhap/src/views/stats"
)

type tab int

const (
	statsTab = iota
	infoTab
	errorTab
	poolsTab
	sessionsTab
	certsTab
	threadsTab
)

type model struct {
	table        table.Model
	message      string
	viewport     viewport.Model
	activeTab    tab
	tabs         []string
	info         string
	errors       string
	pools        string
	certs        string
	threads      string
	sessions     string
	err          error
	lastFetch    time.Time
	width        int
	height       int
	config       Config
	showHelp     bool
	filterMode   bool
	filterInput  string
	allStatsRows []table.Row // Store all rows for filtering
}

type (
	infoMsg    string
	errorMsg   string
	poolsMsg   string
	sessionMsg string
	certsMsg   string
	threadsMsg string
)

type clearMessageMsg struct{}

func fetchStats(cfg Config) tea.Msg {
	conn, err := net.Dial("unix", cfg.socketPath)
	if err != nil {
		log.Printf("Failed to connect to HAProxy socket %s: %v", cfg.socketPath, err)
		return err
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "show stat\n")
	if err != nil {
		log.Printf("Failed to write command to HAProxy socket: %v", err)
		return err
	}

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
		if len(fields) < MinStatsFields {
			log.Printf("Warning: stats line has only %d fields, expected at least %d", len(fields), MinStatsFields)
			continue
		}

		row := table.Row{
			fields[0],              // Name
			fields[1],              // Server
			fields[17],             // Status (will be colored by custom renderer)
			fields[4],              // Current Sessions
			fields[5],              // Max Sessions
			fields[7],              // Total Sessions
			formatBytes(fields[8]), // Bytes In
			formatBytes(fields[9]), // Bytes Out
			fields[13],             // Errors (will be colored by custom renderer)
			fields[18],             // Weight
		}

		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from HAProxy socket: %v", err)
		return err
	}

	return rows
}

func main() {
	// Load config from file
	appConfig := LoadConfig()

	cfg := Config{
		socketPath: appConfig.SocketPath,
	}

	// Command-line argument overrides config file
	if len(os.Args) > 1 {
		cfg.socketPath = os.Args[1]
	}

	// Initial state
	m := model{
		table:     stats.InitializeTable(),
		viewport:  viewport.New(DefaultViewportWidth, DefaultViewportHeight),
		tabs:      []string{"Stats", "Info", "Errors", "Memory", "Sessions", "Certs", "Threads"},
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
