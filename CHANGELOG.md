# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-04-14

### Added

- Session limit (maxconn) column in Stats tab
- Color-coded session counts when approaching or hitting maxconn (yellow at 80%, red at 100%)
- Activity and Events tabs
- Sorting in Stats tab (press `s` to cycle columns)
- Server actions: disable, drain, enable, ready, kill sessions, set weight
- Viewport filtering for all plaintext tabs
- GoReleaser workflow for automated release builds

### Changed

- Upgraded to Charmbracelet v2 (bubbletea, bubbles, lipgloss)
- Updated charmbracelet libraries to latest versions

## [0.2.0] - 2025-02-20

### Added

- Color-coded server status (UP/DOWN/MAINT/DRAIN/NOLB)
- Color-coded error counts (yellow < 10, bold red >= 10)
- Type column with icons for frontends, backends, and servers
- Bold styling for aggregate rows (FE/BE)
- Rate/s column showing sessions per second
- Interactive filtering in Stats and Info tabs (press `/`)
- Vim-style navigation (j/k keys)
- Quick tab jumping (1-9 number keys)
- Help screen (press `?`)
- Config file support (~/.config/lazyhap/config.json)
- Syntax highlighting for Sessions, Threads, and Memory tabs
- SSL Certs and Threads tabs

### Changed

- Reorganized codebase into view packages

### Fixed

- Trimmed value column spaces in Info table

## [0.1.0] - 2025-02-20

### Added

- Initial release
