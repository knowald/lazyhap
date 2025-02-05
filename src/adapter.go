package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
