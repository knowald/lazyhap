package colorize

import (
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
)

// ColorizeSessionOutput adds syntax highlighting to session output
func ColorizeSessionOutput(content string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))        // Cyan for keys
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))      // Green for values
	addrStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))       // Magenta for addresses
	stateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))      // Yellow for states

	// Patterns
	addrPattern := regexp.MustCompile(`0x[0-9a-f]+`)
	keyValuePattern := regexp.MustCompile(`(\w+)=([^\s,\[\]]+)`)

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			result.WriteString(line)
		} else {
			coloredLine := line

			// Colorize addresses (0x...)
			coloredLine = addrPattern.ReplaceAllStringFunc(coloredLine, func(addr string) string {
				return addrStyle.Render(addr)
			})

			// Colorize key=value pairs
			coloredLine = keyValuePattern.ReplaceAllStringFunc(coloredLine, func(match string) string {
				parts := strings.SplitN(match, "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]

					// Special handling for state-like values
					switch value {
					case "EST", "IDLE", "READY", "DONE":
						return keyStyle.Render(key) + "=" + stateStyle.Render(value)
					case "<NONE>", "<none>", "(nil)":
						return keyStyle.Render(key) + "=" + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(value)
					default:
						return keyStyle.Render(key) + "=" + valueStyle.Render(value)
					}
				}
				return match
			})

			result.WriteString(coloredLine)
		}

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// ColorizeThreadOutput adds syntax highlighting to thread output
func ColorizeThreadOutput(content string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))        // Cyan for keys
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))      // Green for values
	addrStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))       // Magenta for addresses
	threadStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true) // Yellow bold for thread headers

	// Patterns
	threadPattern := regexp.MustCompile(`^(\*?\s*Thread\s+\d+\s*:)`)
	addrPattern := regexp.MustCompile(`0x[0-9a-f]+`)
	keyValuePattern := regexp.MustCompile(`(\w+)=([^\s,\[\]()]+)`)

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			result.WriteString(line)
		} else {
			coloredLine := line

			// Highlight thread headers
			if threadPattern.MatchString(line) {
				coloredLine = threadPattern.ReplaceAllStringFunc(coloredLine, func(header string) string {
					return threadStyle.Render(header)
				})
			}

			// Colorize addresses (0x...)
			coloredLine = addrPattern.ReplaceAllStringFunc(coloredLine, func(addr string) string {
				return addrStyle.Render(addr)
			})

			// Colorize key=value pairs
			coloredLine = keyValuePattern.ReplaceAllStringFunc(coloredLine, func(match string) string {
				parts := strings.SplitN(match, "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]

					// Special handling for specific values
					switch value {
					case "(nil)", "<none>", "<NONE>":
						return keyStyle.Render(key) + "=" + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(value)
					case "0", "-1":
						return keyStyle.Render(key) + "=" + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(value)
					default:
						return keyStyle.Render(key) + "=" + valueStyle.Render(value)
					}
				}
				return match
			})

			result.WriteString(coloredLine)
		}

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
