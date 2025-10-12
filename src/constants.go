package main

import "time"

const (
	// Default socket path for HAProxy
	DefaultSocketPath = "/var/run/haproxy/admin.sock"

	// UI dimensions
	DefaultTableHeight    = 20
	DefaultViewportWidth  = 80
	DefaultViewportHeight = 20

	// Timing
	RefreshInterval      = 5 * time.Second
	MessageDisplayTime   = 2 * time.Second
	RetryConnectionDelay = 5 * time.Second

	// HAProxy stats format
	MinStatsFields = 80

	// Server weight
	DefaultServerWeight = 100
)
