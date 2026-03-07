package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	tea "charm.land/bubbletea/v2"
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

func fetchCerts(cfg Config) tea.Msg {
	return certsMsg(execCommand(cfg, "show ssl cert"))
}

func fetchThreads(cfg Config) tea.Msg {
	return threadsMsg(execCommand(cfg, "show threads"))
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

func execCommand(cfg Config, cmd string) string {
	conn, err := net.Dial("unix", cfg.socketPath)
	if err != nil {
		log.Printf("Failed to connect to HAProxy socket %s: %v", cfg.socketPath, err)
		return fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		log.Printf("Failed to write command to HAProxy socket: %v", err)
		return fmt.Sprintf("Error: %v", err)
	}

	var result strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		result.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from HAProxy socket: %v", err)
		return fmt.Sprintf("Error: %v", err)
	}

	return result.String()
}
