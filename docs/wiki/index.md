# arbol Wiki

Welcome. This wiki follows the Jules Dev Standard and FMG Development Bible.

**Current version**: 0.6.0

## Table of Contents

- [Wiki Index (Home)](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/index.md)
- [Architecture & ADRs](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/architecture.md)
- [Development Guide](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/development.md)
- [Agent SOP](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/agent-sop.md)
- [Hygiene & Git Workflow](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/hygiene.md)
- [Roadmap](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/roadmap.md)

## What's New in 0.6.0

- **Code Health** — removed unused code: `noFrame` parameter, `getTerminalWidth`, `padString`, and stale `DIM` variable
- **Security Hardening** — ANSI escape stripping now handles all CSI terminators (0x40-0x7E) instead of only `m`
- **DoS Prevention** — `getDisk` and `getGPU` commands have a 2-second timeout to prevent hanging on unresponsive hardware
- **Performance Boost** — `getProcesses` on Linux uses `syscall.Sysinfo` (~83x faster, from 92µs to 1µs)
- **Default Bar Style** — changed from Braille to Block for better cross-terminal readability
- **Release Automation** — goreleaser workflow builds cross-platform binaries (linux/darwin, amd64/arm64)
- **`--version` flag** — `arbol --version` now reports the installed version
