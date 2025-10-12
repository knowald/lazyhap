package main

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected string
	}{
		{
			name:     "shorter than length",
			input:    "hello",
			length:   10,
			expected: "hello",
		},
		{
			name:     "equal to length",
			input:    "hello",
			length:   5,
			expected: "hello",
		},
		{
			name:     "longer than length",
			input:    "hello world",
			length:   8,
			expected: "hello...",
		},
		{
			name:     "empty string",
			input:    "",
			length:   5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q; want %q", tt.input, tt.length, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "zero bytes",
			input:    "0",
			expected: "0 B",
		},
		{
			name:     "bytes",
			input:    "512",
			expected: "512 B",
		},
		{
			name:     "kilobytes",
			input:    "1024",
			expected: "1.0 KB",
		},
		{
			name:     "megabytes",
			input:    "1048576",
			expected: "1.0 MB",
		},
		{
			name:     "gigabytes",
			input:    "1073741824",
			expected: "1.0 GB",
		},
		{
			name:     "mixed KB",
			input:    "1536",
			expected: "1.5 KB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.input)
			if result != tt.expected {
				t.Errorf("formatBytes(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "positive number",
			input:    "12345",
			expected: 12345,
		},
		{
			name:     "negative number",
			input:    "-100",
			expected: -100,
		},
		{
			name:     "invalid input",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringToInt(tt.input)
			if result != tt.expected {
				t.Errorf("stringToInt(%q) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}
