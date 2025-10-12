package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// AppConfig represents the application configuration
type AppConfig struct {
	SocketPath      string        `json:"socket_path"`
	RefreshInterval time.Duration `json:"refresh_interval_ms"` // in milliseconds
}

// DefaultConfig returns the default configuration
func DefaultConfig() AppConfig {
	return AppConfig{
		SocketPath:      DefaultSocketPath,
		RefreshInterval: RefreshInterval,
	}
}

// LoadConfig loads configuration from file, falling back to defaults
func LoadConfig() AppConfig {
	config := DefaultConfig()

	// Try to load from config file
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file doesn't exist or can't be read, use defaults
		return config
	}

	// Parse JSON config
	var fileConfig struct {
		SocketPath         string `json:"socket_path"`
		RefreshIntervalMs  int    `json:"refresh_interval_ms"`
	}

	if err := json.Unmarshal(data, &fileConfig); err != nil {
		// Invalid config file, use defaults
		return config
	}

	// Apply config from file
	if fileConfig.SocketPath != "" {
		config.SocketPath = fileConfig.SocketPath
	}
	if fileConfig.RefreshIntervalMs > 0 {
		config.RefreshInterval = time.Duration(fileConfig.RefreshIntervalMs) * time.Millisecond
	}

	return config
}

// SaveConfig saves the current configuration to file
func SaveConfig(config AppConfig) error {
	configPath := getConfigPath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Convert to JSON-friendly format
	fileConfig := struct {
		SocketPath        string `json:"socket_path"`
		RefreshIntervalMs int    `json:"refresh_interval_ms"`
	}{
		SocketPath:        config.SocketPath,
		RefreshIntervalMs: int(config.RefreshInterval / time.Millisecond),
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(fileConfig, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	// Try XDG_CONFIG_HOME first
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, "lazyhap", "config.json")
	}

	// Fall back to ~/.config
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", "lazyhap", "config.json")
	}

	// Last resort: current directory
	return "lazyhap.json"
}
