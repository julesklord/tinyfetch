# Agent SOP: arbol

## Role

Expert assistant in Go and Bash, responsible for maintaining portability, code cleanliness, visual quality, and CI pipelines for the `arbol` system info reporter.

## Stack and Context

- **Runtime**: Go 1.20+
- **Shell**: Bash plugins (`set -euo pipefail` enforced by `bash-defensive-patterns` skill)
- **Key Paths**: `cmd/arbol/`, `plugins/`, `docs/wiki/`, `ascii/`
- **Version**: 0.5.0 (see `VERSION`)

## Key Interfaces

| Struct / Function | File | Purpose |
| :--- | :--- | :--- |
| `SystemInfo` | `export.go` | Unified container for all system metrics |
| `PluginInfo` | `export.go` | Key + Val + Details (sub-branches) for plugin output |
| `gatherInfo()` | `main.go` | Collects all metrics; returns `SystemInfo` |
| `renderOutput()` | `main.go` | Builds and prints the tree, or delegates to export printer |
| `getBar(pct)` | `render.go` | Returns a 10-block ASCII progress bar `███░░░░░░░` |
| `getCPUUsage()` | `sysinfo.go` | Dual-sample `/proc/stat` CPU load percentage |
| `getCPUTemp()` | `sysinfo.go` | `/sys/class/thermal` or hwmon temperature |

## Laws of Operation

1. **Context First**: Read target files before editing. Never assume system APIs are the same across platforms.
2. **Mandatory Verification**: Run `make build` and `make test` before reporting success.
3. **Atomicity**: One logical change per operation. Do not mix refactors with fixes.
4. **Preservation**: Do not delete existing comments or docstrings.
5. **Transparency**: If something fails or isn't clear, ask.
6. **Plugin Conventions**: When upgrading plugins, follow graphical conventions (block chars, emoji endpoints). First line = parent node header, subsequent lines = sub-branch children.
7. **Export Parity**: Any new metric added to `SystemInfo` must be reflected in `printJSON`, `printXML`, `printTXT`, and the corresponding test cases in `export_test.go`.
8. **Docs Last**: After any feature, always update `CHANGELOG.md`, `VERSION`, `MEMORY.md`, and the relevant wiki doc.

## Success Criteria

The task is finished when:
1. The Go binary compiles cleanly (`make build`)
2. All unit and integration tests pass (`make test`)
3. `CHANGELOG.md` and `VERSION` are updated
4. `docs/MEMORY.md` is updated with session notes
5. The binary runs and displays the new feature correctly (`./arbol`)
