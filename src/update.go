package main

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

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
			m.activeTab = (m.activeTab + 1) % 5
			if m.activeTab == infoTab {
				m.table = initializeInfoTable()
				rows := parseInfoToRows(m.info)
				m.table.SetRows(rows)
				m.viewport.SetContent(m.info)
				return m, nil
			} else if m.activeTab == statsTab {
				oldRows := m.table.Rows()
				m.table = initializeStatsTable()
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
			m.activeTab = (m.activeTab - 1 + 5) % 5
			if m.activeTab == infoTab {
				m.table = initializeInfoTable()
				rows := parseInfoToRows(m.info)
				m.table.SetRows(rows)
				m.viewport.SetContent(m.info)
				return m, nil
			} else if m.activeTab == statsTab {
				oldRows := m.table.Rows()
				m.table = initializeStatsTable()
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
