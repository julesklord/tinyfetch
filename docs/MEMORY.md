# Project Memory

## Context Index

This file serves as persistent memory for long sessions with AI agents. Don't forget where you came from.

## Current Status

- **Version**: 0.5.0
- **Current Milestone**: 0.5.0 — Graphical plugin metrics, CPU resources, and project renaming complete.
- **Blockers**: None.
- **Next Step**: Polish structured export parity (cpu_usage/cpu_temp in all export formats), consider `--config` support for layout preferences.

## Architecture Quick Reference

- **Binary**: `./arbol` (compiled from `cmd/arbol/`)
- **Plugins dir**: `./plugins/` (simple) and `./plugins/extended/` (diagnostic pane)
- **Plugin discovery**: Env var `ARBOL_PLUGINS_DIR` or sibling `plugins/` directory
- **Banner assets**: Searched in `./ascii/`, `~/.local/share/arbol/ascii/`, `/usr/local/share/arbol/ascii/`, `/usr/share/arbol/ascii/`
- **CPU usage**: Dual-sample `/proc/stat` with 50ms delta (Linux); `ps -A -o %cpu` aggregate (macOS)
- **CPU temp**: `/sys/class/thermal/thermal_zone*/temp`, then `/sys/class/hwmon/hwmon*/temp*_input`
- **Package variable**: `ARBOL_PLUGINS_DIR`

## Session Notes

- 2026-06-10: Initial project analysis and general quality review.
- 2026-06-10: Refactored `scripts/arbol.sh` for multi-platform robustness (Linux & macOS/Darwin) under `set -e`, fixed the ShellCheck SC2034 warning, and added a side-by-side colorized system logo.
- 2026-06-10: Initialized FMG standard files: `docs/AGENT.md`, `docs/GEMINI.md`, `docs/IDENTITY.md`, `docs/MEMORY.md`, and `docs/SOUL.md`.
- 2026-06-11: Implemented Go binary, concurrent plugin loader, tree renderer, unit tests, and CI workflow. Reached 0.4.0.
- 2026-06-28: Implemented multi-line plugin subtrees (parent/child stdout parsing), upgraded all 8 plugins to rich multi-line diagnostic outputs.
- 2026-06-28: Designed flat solid-block TrueColor ASCII art banners for CachyOS, Ubuntu, Arch, Debian, Fedora, openSUSE, Manjaro.
- 2026-06-28: Renamed entire project from `tinyfetch` to `arbol` (module path, binary, env vars, XML tag, docs, Makefile, tests).
- 2026-06-28: Added CPU usage (dual-sample /proc/stat), CPU temperature (/sys/class/thermal), network Rx/Tx counters, packages ratio bar, and weather thermometer scale. Version bumped to 0.5.0.
