package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/knowald/lazyhap/src/views/info"
	"github.com/knowald/lazyhap/src/views/stats"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return fetchStats(m.config) },
		func() tea.Msg { return fetchInfo(m.config) },
		func() tea.Msg { return fetchErrors(m.config) },
		func() tea.Msg { return fetchPools(m.config) },
		func() tea.Msg { return fetchSessions(m.config) },
		func() tea.Msg { return fetchCerts(m.config) },
		func() tea.Msg { return fetchThreads(m.config) },
		func() tea.Msg { return fetchActivity(m.config) },
		func() tea.Msg { return fetchEvents(m.config) },
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
		m.viewport.SetHeight(m.height - headerHeight - footerHeight)
		m.viewport.SetWidth(m.width - 4)

		m.table.SetHeight(m.height - headerHeight - footerHeight)
		m.table.SetWidth(m.width - 4)

	case error:
		m.err = msg
		m.connected = false
		return m, tea.Tick(RetryConnectionDelay, func(t time.Time) tea.Msg {
			return fetchStats(m.config)
		})

	case []table.Row:
		m.connected = true
		m.err = nil
		if m.activeTab == statsTab {
			m.lastFetch = time.Now()
			m.allStatsRows = msg
			m.applySortAndFilter()
		}
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchStats(m.config)
		})

	case infoMsg:
		m.info = string(msg)
		m.allInfoRows = info.ParseInfoToRows(m.info)
		if m.activeTab == infoTab {
			if m.filterMode && m.filterInput != "" {
				m.table.SetRows(filterRows(m.allInfoRows, m.filterInput))
			} else {
				m.table.SetRows(m.allInfoRows)
			}
		}
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchInfo(m.config)
		})

	case errorMsg:
		m.errors = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchErrors(m.config)
		})

	case poolsMsg:
		m.pools = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchPools(m.config)
		})

	case sessionMsg:
		m.sessions = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchSessions(m.config)
		})

	case certsMsg:
		m.certs = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchCerts(m.config)
		})

	case threadsMsg:
		m.threads = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchThreads(m.config)
		})

	case activityMsg:
		m.activity = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchActivity(m.config)
		})

	case eventsMsg:
		m.events = string(msg)
		return m, tea.Tick(RefreshInterval, func(t time.Time) tea.Msg {
			return fetchEvents(m.config)
		})

	case clearMessageMsg:
		m.message = ""
		return m, nil

	case tea.KeyPressMsg:
		// Handle confirm mode
		if m.confirmMode {
			switch msg.String() {
			case "y":
				m.confirmMode = false
				if m.confirmAction == "kill" {
					return m, killServerSessions(m.config, m.confirmBackend, m.confirmServer)
				}
			case "n", "esc":
				m.confirmMode = false
			}
			return m, nil
		}

		// Handle weight input mode
		if m.weightMode {
			switch msg.String() {
			case "enter":
				m.weightMode = false
				if m.weightInput != "" {
					w, err := strconv.Atoi(m.weightInput)
					if err != nil || w < 0 || w > 256 {
						m.message = "Invalid weight (must be 0-256)"
						return m, tea.Tick(MessageDisplayTime, func(t time.Time) tea.Msg {
							return clearMessageMsg{}
						})
					}
					return m, setServerWeight(m.config, m.weightBackend, m.weightServer, w)
				}
				return m, nil
			case "esc":
				m.weightMode = false
				m.weightInput = ""
				return m, nil
			case "backspace":
				if len(m.weightInput) > 0 {
					m.weightInput = m.weightInput[:len(m.weightInput)-1]
				}
				return m, nil
			default:
				c := msg.String()
				if len(c) == 1 && c[0] >= '0' && c[0] <= '9' && len(m.weightInput) < 3 {
					m.weightInput += c
				}
				return m, nil
			}
		}

		// Handle filter input
		if m.filterMode && (m.activeTab == statsTab || m.activeTab == infoTab) {
			switch msg.String() {
			case "enter":
				m.filterMode = false
				return m, nil
			case "esc":
				m.filterMode = false
				m.filterInput = ""
				if m.activeTab == statsTab {
					m.table.SetRows(m.allStatsRows)
				} else {
					m.table.SetRows(m.allInfoRows)
				}
				return m, nil
			case "backspace":
				if len(m.filterInput) > 0 {
					m.filterInput = m.filterInput[:len(m.filterInput)-1]
					m.applyFilter()
				}
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.filterInput += msg.String()
					m.applyFilter()
				}
				return m, nil
			}
		}

		// Handle viewport filter input
		if m.viewportFilterMode {
			switch msg.String() {
			case "enter":
				m.viewportFilterMode = false
				return m, nil
			case "esc":
				m.viewportFilterMode = false
				m.viewportFilterInput = ""
				m.applyViewportFilter()
				return m, nil
			case "backspace":
				if len(m.viewportFilterInput) > 0 {
					m.viewportFilterInput = m.viewportFilterInput[:len(m.viewportFilterInput)-1]
					m.applyViewportFilter()
				}
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.viewportFilterInput += msg.String()
					m.applyViewportFilter()
				}
				return m, nil
			}
		}

		switch msg.String() {
		case "/":
			if m.activeTab == statsTab || m.activeTab == infoTab {
				m.filterMode = true
				m.filterInput = ""
				return m, nil
			}
			// Viewport tabs support filtering too
			m.viewportFilterMode = true
			m.viewportFilterInput = ""
			return m, nil
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "q", "ctrl+c", "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			return m, tea.Quit
		case "j", "down":
			// Forward to table/viewport for navigation
		case "k", "up":
			// Forward to table/viewport for navigation
		case "d":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					backend := selectedRow[1]
					server := selectedRow[2]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, disableServer(m.config, backend, server)
					}
				}
			}
		case "D":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					backend := selectedRow[1]
					server := selectedRow[2]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, drainServer(m.config, backend, server)
					}
				}
			}
		case "e":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					backend := selectedRow[1]
					server := selectedRow[2]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, enableServer(m.config, backend, server)
					}
				}
			}
		case "R":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					backend := selectedRow[1]
					server := selectedRow[2]
					if server != "FRONTEND" && server != "BACKEND" {
						return m, readyServer(m.config, backend, server)
					}
				}
			}
		case "x":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					backend := selectedRow[1]
					server := selectedRow[2]
					if server != "FRONTEND" && server != "BACKEND" {
						m.confirmMode = true
						m.confirmAction = "kill"
						m.confirmBackend = backend
						m.confirmServer = server
						return m, nil
					}
				}
			}
		case "c":
			if m.activeTab == statsTab {
				m.message = "Counters cleared"
				return m, tea.Batch(
					clearCounters(m.config),
					tea.Tick(MessageDisplayTime, func(t time.Time) tea.Msg {
						return clearMessageMsg{}
					}),
				)
			}
		case "w":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					backend := selectedRow[1]
					server := selectedRow[2]
					if server != "FRONTEND" && server != "BACKEND" {
						m.weightMode = true
						m.weightInput = ""
						m.weightBackend = backend
						m.weightServer = server
						return m, nil
					}
				}
			}
		case "s":
			if m.activeTab == statsTab {
				numCols := len(m.table.Columns())
				if m.sortColumn == -1 {
					m.sortColumn = 0
					m.sortAscending = true
				} else if m.sortAscending {
					m.sortAscending = false
				} else {
					m.sortColumn++
					m.sortAscending = true
					if m.sortColumn >= numCols {
						m.sortColumn = -1
					}
				}
				m.applySortAndFilter()
				return m, nil
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Quick jump to tab by number
			tabNum := int(msg.String()[0] - '1')
			if tabNum >= 0 && tabNum < len(m.tabs) {
				previousTab := m.activeTab
				m.activeTab = tab(tabNum)
				m.filterMode = false
				m.filterInput = ""

				// Handle tab switching logic
				if m.activeTab == infoTab {
					m.table = info.InitializeTable()
					m.applyTableSize()
					m.table.SetRows(m.allInfoRows)
					return m, nil
				} else if m.activeTab == statsTab {
					oldRows := m.table.Rows()
					m.table = stats.InitializeTable()
					m.applyTableSize()
					if len(oldRows) > 0 {
						m.table.SetRows(oldRows)
					}
					if previousTab != statsTab {
						return m, func() tea.Msg {
							return fetchStats(m.config)
						}
					}
				}
				return m, nil
			}
		case "tab", "right", "l":
			previousTab := m.activeTab
			m.activeTab = tab((int(m.activeTab) + 1) % len(m.tabs))
			m.filterMode = false
			m.filterInput = ""
			m.viewportFilterMode = false
			m.viewportFilterInput = ""
			if m.activeTab == infoTab {
				m.table = info.InitializeTable()
				m.applyTableSize()
				m.table.SetRows(m.allInfoRows)
				return m, nil
			} else if m.activeTab == statsTab {
				oldRows := m.table.Rows()
				m.table = stats.InitializeTable()
				m.applyTableSize()
				if len(oldRows) > 0 {
					m.table.SetRows(oldRows)
				}
				if previousTab != statsTab {
					return m, func() tea.Msg {
						return fetchStats(m.config)
					}
				}
			}
			return m, nil
		case "shift+tab", "left", "h":
			previousTab := m.activeTab
			m.activeTab = tab((int(m.activeTab) - 1 + len(m.tabs)) % len(m.tabs))
			m.filterMode = false
			m.filterInput = ""
			m.viewportFilterMode = false
			m.viewportFilterInput = ""
			if m.activeTab == infoTab {
				m.table = info.InitializeTable()
				m.applyTableSize()
				m.table.SetRows(m.allInfoRows)
				return m, nil
			} else if m.activeTab == statsTab {
				oldRows := m.table.Rows()
				m.table = stats.InitializeTable()
				m.applyTableSize()
				if len(oldRows) > 0 {
					m.table.SetRows(oldRows)
				}
				if previousTab != statsTab {
					return m, func() tea.Msg {
						return fetchStats(m.config)
					}
				}
				return m, nil
			}
		case "r":
			switch m.activeTab {
			case statsTab:
				return m, func() tea.Msg { return fetchStats(m.config) }
			case infoTab:
				return m, func() tea.Msg { return fetchInfo(m.config) }
			case errorTab:
				return m, func() tea.Msg { return fetchErrors(m.config) }
			case poolsTab:
				return m, func() tea.Msg { return fetchPools(m.config) }
			case sessionsTab:
				return m, func() tea.Msg { return fetchSessions(m.config) }
			case certsTab:
				return m, func() tea.Msg { return fetchCerts(m.config) }
			case threadsTab:
				return m, func() tea.Msg { return fetchThreads(m.config) }
			case activityTab:
				return m, func() tea.Msg { return fetchActivity(m.config) }
			case eventsTab:
				return m, func() tea.Msg { return fetchEvents(m.config) }
			}
		case "g":
			if m.activeTab == statsTab || m.activeTab == infoTab {
				m.table.GotoTop()
			} else {
				m.viewport.GotoTop()
			}
			return m, nil
		case "G":
			if m.activeTab == statsTab || m.activeTab == infoTab {
				m.table.GotoBottom()
			} else {
				m.viewport.GotoBottom()
			}
			return m, nil
		case "y":
			if m.activeTab == statsTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					value := selectedRow[1] + "/" + selectedRow[2]
					err := copyToClipboard(value)
					if err != nil {
						m.err = err
					} else {
						m.message = "Copied to clipboard"
						return m, tea.Tick(MessageDisplayTime, func(t time.Time) tea.Msg {
							return clearMessageMsg{}
						})
					}
				}
			} else if m.activeTab == infoTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 2 {
					value := selectedRow[1]
					err := copyToClipboard(value)
					if err != nil {
						m.err = err
					} else {
						m.message = "Copied to clipboard"
						return m, tea.Tick(MessageDisplayTime, func(t time.Time) tea.Msg {
							return clearMessageMsg{}
						})
					}
				}
			}
		}
	}

	switch m.activeTab {
	case statsTab, infoTab:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) applyTableSize() {
	headerHeight := 4
	footerHeight := 2
	m.table.SetHeight(m.height - headerHeight - footerHeight)
	m.table.SetWidth(m.width - 4)
}

func (m *model) applyFilter() {
	if m.activeTab == statsTab {
		m.applySortAndFilter()
	} else if m.activeTab == infoTab {
		m.table.SetRows(filterRows(m.allInfoRows, m.filterInput))
	}
}

func (m *model) applyViewportFilter() {
	var content string
	switch m.activeTab {
	case errorTab:
		content = m.errors
	case sessionsTab:
		content = m.sessions
	case certsTab:
		content = m.certs
	case threadsTab:
		content = m.threads
	case activityTab:
		content = m.activity
	case eventsTab:
		content = m.events
	case poolsTab:
		content = m.pools
	}
	if m.viewportFilterInput != "" {
		content = filterViewportLines(content, m.viewportFilterInput)
	}
	m.viewport.SetContent(content)
}

func (m *model) applySortAndFilter() {
	rows := m.allStatsRows
	if m.sortColumn >= 0 {
		rows = sortRows(rows, m.sortColumn, m.sortAscending)
	}
	if m.filterMode && m.filterInput != "" {
		rows = filterRows(rows, m.filterInput)
	}
	m.table.SetRows(rows)
}

func sortRows(rows []table.Row, col int, ascending bool) []table.Row {
	sorted := make([]table.Row, len(rows))
	copy(sorted, rows)

	sort.SliceStable(sorted, func(i, j int) bool {
		a := sorted[i][col]
		b := sorted[j][col]

		// Strip ANSI for comparison
		a = stripANSI(strings.TrimSpace(a))
		b = stripANSI(strings.TrimSpace(b))

		// Try numeric comparison first
		aNum, aErr := strconv.ParseFloat(a, 64)
		bNum, bErr := strconv.ParseFloat(b, 64)

		// Handle byte-formatted values (e.g. "1.2 KB")
		if aErr != nil {
			aNum, aErr = parseByteValue(a)
		}
		if bErr != nil {
			bNum, bErr = parseByteValue(b)
		}

		if aErr == nil && bErr == nil {
			if ascending {
				return aNum < bNum
			}
			return aNum > bNum
		}

		// Fall back to string comparison
		if ascending {
			return strings.ToLower(a) < strings.ToLower(b)
		}
		return strings.ToLower(a) > strings.ToLower(b)
	})

	return sorted
}

func parseByteValue(s string) (float64, error) {
	s = strings.TrimSpace(s)
	multipliers := map[string]float64{
		"B": 1, "KB": 1024, "MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024, "TB": 1024 * 1024 * 1024 * 1024,
	}
	for suffix, mult := range multipliers {
		if strings.HasSuffix(s, " "+suffix) {
			numStr := strings.TrimSuffix(s, " "+suffix)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			return num * mult, nil
		}
	}
	return 0, strconv.ErrSyntax
}

// Model get/set

func (m model) GetViewport() viewport.Model {
	return m.viewport
}

func (m model) ErrorView() string {
	return m.errors
}

func (m model) TableView() string {
	return m.table.View()
}

func (m model) LastFetchTime() string {
	return m.lastFetch.Format("15:04:05")
}

func (m model) GetMessage() string {
	return m.message
}

func (m model) CertsView() string {
	return m.certs
}

func (m model) ThreadsView() string {
	return m.threads
}

func (m model) PoolsView() string {
	return m.pools
}

func (m model) SessionsView() string {
	return m.sessions
}

func (m model) ActivityView() string {
	return m.activity
}

func (m model) EventsView() string {
	return m.events
}

func (m model) FilterMode() bool {
	return m.filterMode
}

func (m model) FilterInput() string {
	return m.filterInput
}

func (m model) GetTable() table.Model {
	return m.table
}

func (m model) ConfirmMode() bool {
	return m.confirmMode
}

func (m model) ConfirmPrompt() string {
	return "Kill all sessions on " + m.confirmBackend + "/" + m.confirmServer + "? (y/n)"
}

func (m model) WeightMode() bool {
	return m.weightMode
}

func (m model) WeightInput() string {
	return m.weightInput
}

func (m model) WeightServer() string {
	return m.weightBackend + "/" + m.weightServer
}

func (m model) ViewportFilterMode() bool {
	return m.viewportFilterMode
}

func (m model) ViewportFilterInput() string {
	return m.viewportFilterInput
}

func (m model) Connected() bool {
	return m.connected
}

func (m model) SocketPath() string {
	return m.config.socketPath
}

func (m model) SortColumn() int {
	return m.sortColumn
}

func (m model) SortAscending() bool {
	return m.sortAscending
}
