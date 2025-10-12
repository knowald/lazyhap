# LazyHAP

## Overview

LazyHAP is a lightweight, portable HAProxy TUI tool written in Go.

![Screenshot](screenshot.png)

## Features

- Real-time HAProxy server statistics with auto-refresh
- **Color-coded server status** (UP/DOWN/MAINT/DRAIN/NOLB)
- **Interactive filtering** in Stats tab (press `/` to search)
- **Quick navigation** with vim-style keys and number shortcuts
- **Built-in help system** (press `?`)
- **Config file support** for persistent preferences
- Multiple tab views:
  - Stats (with server control and filtering)
  - Info (with clipboard copy support)
  - Errors
  - Memory
  - Sessions
  - Certs
  - Threads

## Status

🚧 Early Prototype 🚧

- Experimental implementation
- Subject to significant changes

## Build

```bash
go build
```

## Usage

Use default HAProxy socket path (`/var/run/haproxy/admin.sock`)

```bash
./lazyhap
```

Specify custom socket path

```bash
./lazyhap /path/to/custom/haproxy/admin.sock
```

### Configuration File

LazyHAP supports optional configuration via `~/.config/lazyhap/config.json`:

```json
{
  "socket_path": "/var/run/haproxy/admin.sock",
  "refresh_interval_ms": 5000
}
```

Command-line arguments override config file settings.

## Controls

### Navigation
- `tab`/`shift+tab` or `left`/`right`/`h`/`l`: Navigate tabs
- `1-7`: Quick jump to tab by number
- `j`/`k` or `up`/`down`: Navigate within tables
- `?`: Toggle help screen
- `q`/`esc`/`ctrl+c`: Quit

### Stats Tab Commands
- `/`: Start filtering (type to search, Enter to apply, Esc to clear)
- `d`: Disable selected server
- `e`: Enable selected server
- `w`: Set server weight to 100

### Info Tab Commands
- `y`: Yank (copy) selected value to clipboard

### Visual Indicators
- **Green** status: Server is UP
- **Red** status: Server is DOWN
- **Yellow** status: Server in MAINT mode
- **Cyan** status: Server is DRAIN
- **Magenta** status: Server is NOLB
- **Yellow** errors: Low error count (< 10)
- **Bold Red** errors: High error count (≥ 10)

## Requirements

- HAProxy
- Unix-like system with socket access

## Disclaimer

This is an early-stage project and should not be actually used
anywhere near a production server.

## License

[MIT License](LICENSE)

## Contributions

Contributions and feedback welcome.
