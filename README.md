# arbol

> Terminal system info reporter — a beautiful, plugin-driven utility that shows system stats as a live tree in your terminal.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Built With](https://img.shields.io/badge/Built%20With-Go-brightgreen.svg)](go.mod)
[![Version](https://img.shields.io/badge/version-0.5.0-orange.svg)](VERSION)

## Demo

![arbol demo](docs/demo.gif)

## Overview

**arbol** is a dependency-free CLI status utility designed to quickly fetch and display essential system information in your terminal. It is implemented as a high-performance compiled Go binary with an extensible plugin system.

It presents a side-by-side representation of the host OS logo as a TrueColor banner and core system resource metrics (Host, OS, Kernel, Uptime, Shell, CPU, GPU, DE/WM, Terminal, CPU Usage, CPU Temperature, Memory, Swap, Disk, and Process Count) structured as a dynamic tree. Below the tree, extended plugins render a third diagnostic pane with rich charts and graphs.

## Installation

### From Source

Ensure Go 1.20+ is installed:

```bash
git clone https://github.com/julesklord/arbol.git
cd arbol
make build
```

### System-Wide Installation

Install the compiled binary into `/usr/local/bin`:

```bash
sudo make install
```

## Usage

Run the utility directly from your shell:

```bash
arbol
```

### Command Reference

| Option | Short Alias | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display version and usage instructions. |
| `--no-ascii` | | Omit the system ASCII logo. |
| `--minimal` | | Skip extended plugins and display a single info card. |
| `--noframe` | | Omit the box borders and print layout side-by-side using spaces. |
| `--output=FORMAT` | | Serialize system stats and simple plugins into structured output: `json`, `xml`, or `txt`. |
| `--logo=MODE` | | Control the ASCII logo style: `simple` (default glyph block art) or `banner` (solid filled block art). |
| `--plugins-dir=PATH` | | Override plugin discovery directory (also via `ARBOL_PLUGINS_DIR` env var). |

## What It Shows

### 🌿 Tree Layout (Default)

```
● hostname @ DistroName
├── ⚙ specs
│   ├── 📦 kernel: 6.10.x-cachyos
│   ├── ⏱ uptime: 10h 24m
│   ├── 💻 shell: /bin/fish
│   ├── 🧠 cpu: Intel Xeon E5-2630 @ 2.20GHz
│   ├── 🎮 gpu: NVIDIA GTX 1060 6GB
│   ├── 🖥 de/wm: mango
│   └── 📟 terminal: zed
├── 📊 resources
│   ├── 📈 cpu usage: ██░░░░░░░░ 21%
│   ├── 🌡️ cpu temp: 51.0°C
│   ├── 💾 memory: ███░░░░░░░ 36% (39940MB)
│   ├── 🔄 swap: ░░░░░░░░░░ 0% (2MB / 39939MB)
│   ├── 💿 disk: █████░░░░░ /dev/sda2 (56%)
│   └── ⚡ processes: 497
└── 🔌 plugins
    ├── docker: 🐳 idle
    ├── network: 󰖩 connected
    │   ├── Interface: enp8s0
    │   ├── Download (Rx): 📥 3.04 GB
    │   └── Upload (Tx): 📤 586 MB
    ├── packages: 2034 total
    │   └── Ratio: [█████████░] (native/aur/others)
    └── weather: ☀️ +40°C
        └── Temp Scale: ❄️ ━━━━━━━━━━━█━ 🔥
```

### 🔍 Diagnostics Pane (Extended Plugins)

A third column provides rich diagnostic charts:

- **Git Commit Graph** — last 5 commits shown with `git log --oneline --graph`
- **System Dashboard** — load averages and top memory consumers
- **Weather Forecast** — multi-day ASCII weather art from `wttr.in`

## Plugins

Plugins are executable scripts or binaries placed in `./plugins/`. They output one or multiple lines; the first line becomes the tree node header and subsequent lines become nested sub-branches.

### Built-in Plugins

| Plugin | Description |
| :--- | :--- |
| `battery.sh` | Capacity %, status, and ASCII health bar |
| `docker.sh` | Container counts, image/volume totals, and list of active containers |
| `git.sh` | Branch status, staged/modified/untracked counts, last commit, and sync state |
| `ip.sh` | Local IP, public IP, DNS, gateway, interface, and Rx/Tx traffic totals |
| `k8s.sh` | Kubernetes context, namespace, API server, and config path |
| `media.sh` | Now Playing track, artist, album, and media player |
| `packages.sh` | Cross-platform package counts (pacman, AUR, dpkg, rpm, brew, flatpak, snap) with ratio bar |
| `weather.sh` | Location, condition, temperature, and a dynamic thermometer scale graph |

### Extended Plugins (Diagnostics)

| Plugin | Description |
| :--- | :--- |
| `git_graph.sh` | ASCII git log graph of recent commits |
| `sys_dashboard.sh` | System load average and top memory processes |
| `weather_forecast.sh` | Full ASCII weather art from wttr.in |

### Custom Plugins

Any executable script placed in `./plugins/` is automatically discovered at runtime. No registration required:

```bash
#!/usr/bin/env bash
echo "My Plugin: Hello!"
echo "Detail line one"
echo "Detail line two"
```

See the [Development Guide](docs/wiki/development.md) for the full plugin authoring specification.

## Structured Output

Export system info to structured formats with `--output`:

```bash
arbol --output=json
arbol --output=xml
arbol --output=txt
```

JSON output includes `cpu_usage` and `cpu_temp` fields alongside all core system metrics.

## Environment Variables

| Variable | Description |
| :--- | :--- |
| `ARBOL_PLUGINS_DIR` | Override the default plugin discovery directory |
| `ARBOL_ENABLE_WEATHER` | Set to `1` to enable location-based weather fetching (disabled by default for privacy) |

## License

MIT — see [LICENSE](LICENSE).
