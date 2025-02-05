package main

import (
	"net"
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

// TestServer is a mock Unix socket server for testing
type TestServer struct {
	listener net.Listener
	t        *testing.T
	response string
}

func newTestServer(t *testing.T) (*TestServer, error) {
	// Create temp socket file
	socketPath := "/tmp/haproxy-test.sock"
	os.Remove(socketPath) // Clean up any existing socket

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}

	server := &TestServer{
		listener: listener,
		t:        t,
	}

	// Handle connections in background
	go server.serve()

	return server, nil
}

func (s *TestServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return // listener closed
		}

		go func(c net.Conn) {
			defer c.Close()

			// Read command
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				s.t.Errorf("Failed to read command: %v", err)
				return
			}

			// Send mock response
			_, err = c.Write([]byte(s.response))
			if err != nil {
				s.t.Errorf("Failed to write response: %v", err)
				return
			}
		}(conn)
	}
}

func (s *TestServer) close() {
	s.listener.Close()
	os.Remove(s.listener.Addr().String())
}

func TestFetchStats(t *testing.T) {
	// Start test server
	server, err := newTestServer(t)
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer server.close()

	// Set mock response
	server.response = `# pxname,svname,status,weight
backend1,FRONTEND,OPEN,
backend1,server1,UP,100
backend1,BACKEND,UP,
`

	// Create config pointing to test socket
	cfg := Config{
		socketPath: server.listener.Addr().String(),
	}

	// Call fetchStats
	msg := fetchStats(cfg)

	// Assert response
	rows, ok := msg.([]table.Row)
	if !ok {
		t.Fatalf("Expected []table.Row, got %T", msg)
	}

	// Verify rows
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(rows))
	}

	// Check specific values
	if rows[1][0] != "backend1" || rows[1][1] != "server1" || rows[1][2] != "UP" {
		t.Errorf("Unexpected row content: %v", rows[1])
	}
}
