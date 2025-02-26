package main

import (
	"lazyhap/src/views/info"
	"lazyhap/src/views/stats"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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
		if m.activeTab == statsTab {
			m.lastFetch = time.Now()
			m.table.SetRows(msg)
		}
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchStats(m.config)
		})

	case infoMsg:
		m.info = string(msg)
		if m.activeTab == infoTab {
			rows := parseInfoToRows(m.info)
			m.table.SetRows(rows)
		}
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

	case certsMsg:
		m.certs = string(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchCerts(m.config)
		})

	case threadsMsg:
		m.threads = string(msg)
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchThreads(m.config)
		})

	case clearMessageMsg:
		m.message = ""
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
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
		case "tab", "right", "l":
			previousTab := m.activeTab
			m.activeTab = tab((int(m.activeTab) + 1) % len(m.tabs))
			if m.activeTab == infoTab {
				m.table = info.InitializeTable()
				rows := parseInfoToRows(m.info)
				m.table.SetRows(rows)
				m.viewport.SetContent(m.info)
				return m, nil
			} else if m.activeTab == statsTab {
				oldRows := m.table.Rows()
				m.table = stats.InitializeTable()
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
			if m.activeTab == infoTab {
				m.table = info.InitializeTable()
				rows := parseInfoToRows(m.info)
				m.table.SetRows(rows)
				m.viewport.SetContent(m.info)
				return m, nil
			} else if m.activeTab == statsTab {
				oldRows := m.table.Rows()
				m.table = stats.InitializeTable()
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
		case "y":
			if m.activeTab == infoTab {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 2 {
					value := selectedRow[1]
					err := copyToClipboard(value)
					if err != nil {
						m.err = err
					} else {
						m.message = "✓ Copied to clipboard!"
						return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
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
