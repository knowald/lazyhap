package info

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

func TestParseInfoToRows(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []table.Row
	}{
		{
			name:  "simple key-value",
			input: "Name: HAProxy\nVersion: 2.4.0",
			expected: []table.Row{
				{"Name", "HAProxy"},
				{"Version", "2.4.0"},
			},
		},
		{
			name:  "with description",
			input: "Uptime: 3600: Server uptime in seconds",
			expected: []table.Row{
				{"Uptime", "3600", "Server uptime in seconds"},
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []table.Row{},
		},
		{
			name:  "with blank lines",
			input: "Name: HAProxy\n\nVersion: 2.4.0\n\n",
			expected: []table.Row{
				{"Name", "HAProxy"},
				{"Version", "2.4.0"},
			},
		},
		{
			name:  "with spaces around colons",
			input: "Name : HAProxy \nVersion : 2.4.0",
			expected: []table.Row{
				{"Name", "HAProxy"},
				{"Version", "2.4.0"},
			},
		},
		{
			name:     "missing colon",
			input:    "InvalidLine\nValid: Line",
			expected: []table.Row{{"Valid", "Line"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInfoToRows(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("ParseInfoToRows() returned %d rows; want %d rows", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if len(result[i]) != len(tt.expected[i]) {
					t.Errorf("Row %d has %d columns; want %d columns", i, len(result[i]), len(tt.expected[i]))
					continue
				}
				for j := range result[i] {
					if result[i][j] != tt.expected[i][j] {
						t.Errorf("Row %d, Col %d = %q; want %q", i, j, result[i][j], tt.expected[i][j])
					}
				}
			}
		})
	}
}
