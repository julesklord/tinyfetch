# Development Guide

This document guides developers on local setup, running tests, and creating custom plugins.

## Prerequisites

- Go 1.20+
- A Nerd Font in your terminal emulator (for Nerd Font glyphs in plugins)
- `bash`, `curl`, `ip`/`ifconfig` available for plugin scripts

## Local Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/julesklord/arbol.git
   cd arbol
   ```
2. Build the binary using the Makefile:
   ```bash
   make build
   ```
3. Run the test suite to verify everything compiles and behaves correctly:
   ```bash
   make test
   ```
4. Run it:
   ```bash
   ./arbol
   ```

## Makefile Targets

| Target | Description |
| :--- | :--- |
| `make build` | Compiles the Go binary to `./arbol` |
| `make test` | Builds + runs unit tests (`go test`) and integration tests (`tests/test.sh`) |
| `make install` | Installs binary to `/usr/local/bin/arbol` and assets to `/usr/local/share/arbol/ascii/` |
| `make uninstall` | Removes installed binary and assets |
| `make clean` | Deletes the local `./arbol` binary |

## Project Layout

```
arbol/
├── cmd/arbol/
│   ├── main.go         # Entry point, flag parsing, gatherInfo(), renderOutput()
│   ├── sysinfo.go      # Platform-specific metric collectors (CPU, GPU, memory, disk, etc.)
│   ├── render.go       # Progress bar, ANSI strip, padding helpers
│   ├── export.go       # Structured export: JSON, XML, TXT + SystemInfo/PluginInfo structs
│   └── export_test.go  # Unit tests for all export printers
├── plugins/
│   ├── battery.sh      # Battery capacity, status, health bar
│   ├── docker.sh       # Container counts, image/volume totals, active container list
│   ├── git.sh          # Branch, staged/modified/untracked counts, last commit
│   ├── ip.sh           # Local/public IP, gateway, DNS, Rx/Tx traffic totals
│   ├── k8s.sh          # Kubernetes context, namespace, API server
│   ├── media.sh        # Now playing from playerctl/osascript
│   ├── packages.sh     # Cross-platform package counts with ratio distribution bar
│   ├── weather.sh      # Location, condition, temperature + thermometer scale
│   └── extended/
│       ├── git_graph.sh          # git log --graph framed in a box
│       ├── sys_dashboard.sh      # Load avg + top memory processes
│       └── weather_forecast.sh   # Full ASCII weather art from wttr.in
├── ascii/              # Distro banner text files (e.g., cachyos_banner.txt)
├── docs/
│   ├── wiki/
│   │   ├── architecture.md
│   │   ├── development.md   ← this file
│   │   ├── agent-sop.md
│   │   ├── hygiene.md
│   │   ├── index.md
│   │   └── roadmap.md
│   ├── AGENT.md
│   ├── GEMINI.md
│   ├── IDENTITY.md
│   ├── MEMORY.md
│   └── SOUL.md
├── tests/
│   └── test.sh         # Integration test suite
├── Makefile
├── go.mod
├── VERSION
├── CHANGELOG.md
└── README.md
```

## Creating Custom Plugins

`arbol` scans the `./plugins` directory and its `extended/` subdirectory for executable scripts or binaries. You can write plugins in Bash, Python, Go, Node, or any other scripting language.

### Simple Plugins

1. **Location**: Place your script under `./plugins/` (e.g., `plugins/battery.sh`).
2. **Executability**: The file must be executable. Run `chmod +x plugins/my-plugin` to enable it.
3. **Stdout Format**: The plugin can output a single line or multiple lines.
   - The first line can follow the `Label: Value` pattern. If a colon is detected, the key will be printed in blue, and the value in default colors.
   - If multiple lines are printed, subsequent lines will automatically be parsed and rendered as nested sub-branches under the parent node in the tree.
4. **Error Handling**: If the plugin fails, it must exit silently (`exit 0`) and print nothing. If a plugin prints nothing, the row/node is omitted from the tree.

### Extended Plugins

1. **Location**: Place your script under `./plugins/extended/` (e.g., `plugins/extended/sys_dashboard.sh`).
2. **Executability**: Must be executable (`chmod +x`).
3. **Stdout Format**: Can output multiple lines. `arbol` will dynamically calculate widths and align the borders of the third pane symmetrically.
4. **Error Handling**: If the plugin fails or is not applicable, it must exit silently (`exit 0`) and print nothing. If all extended plugins print nothing, the third column will be cleanly omitted from the output.

### Graphical Conventions for Plugins

Use block characters to create micro-charts inline:

| Character | Use |
| :--- | :--- |
| `█` | Filled segment (native packages, battery, etc.) |
| `▒` | Semi-filled segment (AUR packages, mid-level) |
| `░` | Empty segment (unfilled bar portions) |
| `━` | Horizontal scale track |
| `❄️` / `🔥` | Scale endpoints (cold / hot) |
| `📥` / `📤` | Traffic direction indicators |

### Example Simple Plugin (Shell)

`plugins/my-plugin.sh`:
```bash
#!/usr/bin/env bash
set -euo pipefail

# Print parent node: Label: Value
echo "My Plugin: active"
# Print child sub-branches (any number)
echo "Detail: something useful"
echo "Metric: $(uptime -p)"
```

Make it executable:
```bash
chmod +x plugins/my-plugin.sh
```

### Example Plugin with Graphical Bar

```bash
#!/usr/bin/env bash
set -euo pipefail

pct=72   # replace with your real metric
filled=$(( pct * 10 / 100 ))
bar=""
for ((i=0; i<10; i++)); do
  [ "$i" -lt "$filled" ] && bar="${bar}█" || bar="${bar}░"
done

echo "My Metric: [${bar}] ${pct}%"
```

## Environment Variables

| Variable | Default | Description |
| :--- | :--- | :--- |
| `ARBOL_PLUGINS_DIR` | `./plugins` | Override plugin discovery directory |

## Running Tests

```bash
make test
```

This runs:
1. `go test ./cmd/arbol/...` — unit tests for `export_test.go` and `sysinfo_test.go`
2. `tests/test.sh` — integration tests against the compiled `./arbol` binary

## Git Workflow & Conventions

Follow the rules in [hygiene.md](file:///mnt/DEV/Proyectos/repos/arbol/docs/wiki/hygiene.md) for conventional commit messages. Keep changes atomic and always run `make test` before committing.
