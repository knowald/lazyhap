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
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	statsTab tab = iota
	infoTab
	errorTab
	poolsTab
	sessionsTab
)

var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 1)

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			MarginLeft(2)
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 4
		footerHeight := 2
		m.viewport.Height = m.height - headerHeight - footerHeight
		m.viewport.Width = m.width - 4

		m.table.SetHeight(m.height - headerHeight - footerHeight)

		return m, nil

	case error:
		m.err = msg
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchStats(m.config)
		})

	case []table.Row:
		m.lastFetch = time.Now()
		m.table.SetRows(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchStats(m.config)
		})

	case infoMsg:
		m.info = string(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchInfo(m.config)
		})

	case errorMsg:
		m.errors = string(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchErrors(m.config)
		})

	case poolsMsg:
		m.pools = string(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchPools(m.config)
		})

	case sessionMsg:
		m.sessions = string(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchSessions(m.config)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % 5
			return m, nil
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + 5) % 5
			return m, nil
		case "D", "d":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 2 {
					backend := selectedRow[0]
					server := selectedRow[1]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, disableServer(m.config, backend, server)
					}
				}
			}
		case "e":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 2 {
					backend := selectedRow[0]
					server := selectedRow[1]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, enableServer(m.config, backend, server)
					}
				}
			}
		case "w":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 2 {
					backend := selectedRow[0]
					server := selectedRow[1]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, setServerWeight(m.config, backend, server, 100)
					}
				}
			}
		}
	}

	switch m.activeTab {
	case statsTab:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\nPress q to quit\n", m.err)
	}

	var sb strings.Builder

	renderedTabs := []string{}
	for i, t := range m.tabs {
		if i == int(m.activeTab) {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...))
	sb.WriteString("\n\n")

	switch m.activeTab {
	case statsTab:
		timeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)
		sb.WriteString(timeStyle.Render(fmt.Sprintf("Last updated: %s", m.lastFetch.Format("15:04:05"))))
		sb.WriteString("\n")
		sb.WriteString(baseStyle.Render(m.table.View()))
		sb.WriteString("\n")
		sb.WriteString("Commands: (d)isable server, (e)nable server, set (w)eight to 100")

	case infoTab:
		m.viewport.SetContent(m.info)
		sb.WriteString(baseStyle.Render(m.viewport.View()))

	case errorTab:
		m.viewport.SetContent(m.errors)
		sb.WriteString(baseStyle.Render(m.viewport.View()))

	case poolsTab:
		m.viewport.SetContent(m.pools)
		sb.WriteString(baseStyle.Render(m.viewport.View()))

	case sessionsTab:
		m.viewport.SetContent(m.sessions)
		sb.WriteString(baseStyle.Render(m.viewport.View()))
	}

	return sb.String()
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

		var status string
		if fields[17] == "UP" { // Clean status field
			status = "UP" // styleStatus("UP")
		} else if fields[17] == "DOWN" {
			status = "DOWN" // styleStatus("DOWN")
		} else {
			status = fields[17]
		}

		row := table.Row{
			// truncate(fields[0], 18), // Name
			// truncate(fields[1], 13), // Server
			fields[0],              // Name
			fields[1],              // Server
			status,                 // Status
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
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		rows = append(rows, table.Row{key, value})
	}

	return rows
}

func fetchInfo(cfg Config) tea.Msg {
	return infoMsg(execCommand(cfg, "show info desc"))
}

func fetchErrors(cfg Config) tea.Msg {
	return errorMsg(execCommand(cfg, "show errors"))
}

func fetchPools(cfg Config) tea.Msg {
	return poolsMsg(execCommand(cfg, "show pools"))
}

func fetchSessions(cfg Config) tea.Msg {
	return sessionMsg(execCommand(cfg, "show sess"))
}

func disableServer(cfg Config, backend, server string) tea.Cmd {
	return func() tea.Msg {
		cmd := fmt.Sprintf("disable server %s/%s", backend, server)
		execCommand(cfg, cmd)
		return fetchStats(cfg)
	}
}

func enableServer(cfg Config, backend, server string) tea.Cmd {
	return func() tea.Msg {
		cmd := fmt.Sprintf("enable server %s/%s", backend, server)
		execCommand(cfg, cmd)
		return fetchStats(cfg)
	}
}

func setServerWeight(cfg Config, backend, server string, weight int) tea.Cmd {
	return func() tea.Msg {
		cmd := fmt.Sprintf("set server %s/%s weight %d", backend, server, weight)
		execCommand(cfg, cmd)
		return fetchStats(cfg)
	}
}

type Config struct {
	socketPath string
}

func execCommand(cfg Config, cmd string) string {
	conn, err := net.Dial("unix", cfg.socketPath)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", cmd)

	var result strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		result.WriteString(scanner.Text() + "\n")
	}

	return result.String()
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func styleStatus(status string) string {
	switch strings.ToLower(status) {
	case "up":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Render("UP")
	case "down":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render("DOWN")
	default:
		return status
	}
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

func main() {
	cfg := Config{
		socketPath: "/var/run/haproxy/admin.sock", // default path
	}

	if len(os.Args) > 1 {
		cfg.socketPath = os.Args[1]
	}

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

	t.SetStyles(s)

	m := model{
		table:     t,
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
