# tinyfetch

> Minimal fastfetch-style status tool — a tiny, focused utility to show system info in a compact, beautiful layout.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Built With](https://img.shields.io/badge/Built%20With-Shell%20%2F%20Go-brightgreen.svg)](go.mod)

## Overview

**tinyfetch** is a tiny, dependency-free CLI status utility designed to quickly fetch and display essential system information in your terminal. It offers two implementations:
1. A robust, portable POSIX-compliant Shell script (`scripts/tinyfetch.sh`).
2. A high-performance compiled Go binary.

Both versions present a side-by-side colorized representation of the host OS logo and core resource metrics (Host, OS, Kernel, Uptime, Shell, CPU, Memory, and Disk usage).

## Installation

### Via Script (Quickstart)

```bash
git clone https://github.com/julesklord/tinyfetch.git
cd tinyfetch
scripts/tinyfetch.sh
```

### From Source (Go Version)

Ensure Go 1.20+ is installed:

```bash
make build
```

### System-Wide Installation

Install the script or compiled binary into `/usr/local/bin`:

```bash
sudo make install
```

## Usage

Run the utility directly from your shell:

```bash
tinyfetch
```

### Command Reference

| Option | Short Alias | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display version and usage instructions. |
| `--no-ascii` | | Omit the side-by-side system ASCII logo. |
| `--minimal` | | Skip extended plugins and display a single info card. |
| `--noframe` | | Omit the box borders and print layout side-by-side using spaces. |
| `--output=FORMAT` | | Serialize system stats and simple plugins into structured output: `json`, `xml`, or `txt`. |

## Extensibility & Plugins

**tinyfetch** is fully extensible via custom plugins. It scans the `./plugins` directory for executable scripts or binaries (written in Bash, Python, Go, Node, etc.) and appends their output dynamically to the dashboard.

It supports two types of plugins:
1. **Simple Plugins** (located in `./plugins/`): Single-line status elements that append rows to the middle info card.
2. **Extended Plugins** (located in `./plugins/extended/`): Multi-line, complex dashboards that render side-by-side in a separate third column (creating a 3-pane layout, or a 2-pane layout if `--no-ascii` is passed).

### Creating a Simple Plugin

1. Create an executable file inside the `./plugins/` directory:
   ```bash
   touch plugins/my-plugin.sh
   chmod +x plugins/my-plugin.sh
   ```
2. Your script should output exactly one line in one of the following formats:
   - **Label format**: `Label: Value` (e.g. `Network: Connected`). `tinyfetch` will automatically colorize the label in blue.
   - **Plain format**: `Value` (e.g. `Connected`). `tinyfetch` will automatically format it using the capitalized filename as the label.
3. If a plugin needs to exit early or is not applicable, it should print nothing and exit with code `0`. Any plugin that produces empty output will be cleanly omitted from the dashboard.

### Creating an Extended Plugin

1. Create an executable file inside the `./plugins/extended/` directory:
   ```bash
   touch plugins/extended/my-dashboard.sh
   chmod +x plugins/extended/my-dashboard.sh
   ```
2. Your script can output multiple lines. Both Shell and Go versions of `tinyfetch` will dynamically calculate widths and align the borders of the third pane symmetrically using rune-count calculations.
3. If the script exits with `0` and prints nothing, the third column will be omitted and `tinyfetch` will fall back to the standard double-pane layout automatically.

### Included Plugins

The repository contains several useful out-of-the-box plugins under [plugins/](file:///home/julesklord/Proyectos/repos/mini-fetch/plugins):
- **Battery** (`battery.sh`): Displays current battery percentage and charge status (supports Linux & macOS).
- **Docker** (`docker.sh`): Shows active containers and daemon status.
- **Git** (`git.sh`): Reports current branch, dirty status, and counters for staged/modified/untracked files with Nerd Fonts.
- **Network** (`ip.sh`): Shows local and external IP addresses (using `icanhazip.com` with a 1s connection timeout).
- **Kubernetes** (`k8s.sh`): Reports active kubectl context and namespace.
- **Packages** (`packages.sh`): Lists installed packages (supports `pacman` and `paru`/`yay`).
- **Weather** (`weather.sh`): Fetches temperature and sky status.
- **Media Player** (`media.sh`): Shows currently playing song or media status.

#### Extended Plugins ([plugins/extended/](file:///home/julesklord/Proyectos/repos/mini-fetch/plugins/extended)):
- **Git Commit Graph** (`git_graph.sh`): Displays a beautiful 5-line local branch history tree visualization (or pulls remote status from the GitHub API if `GITHUB_TOKEN` is present).
- **System Dashboard** (`sys_dashboard.sh`): Displays load averages and top memory-consuming processes.
- **Weather Forecast** (`weather_forecast.sh`): Displays a multi-line weather forecast from `wttr.in`.

## Architecture

The utility checks the runtime operating system and dynamically resolves resource usage metrics through the most efficient native queries available.

```mermaid
graph TD
    User([UserActor]) -->|RunCommand| Cli[CliTool]
    Cli -->|QuerySystem| OsCheck{OsCheck}
    OsCheck -->|QueryLinux| LinuxParser[LinuxParser]
    OsCheck -->|QueryDarwin| DarwinParser[DarwinParser]
    LinuxParser -->|ReadFiles| ProcUptime[/proc/uptime]
    LinuxParser -->|ReadFiles| ProcMem[/proc/meminfo]
    LinuxParser -->|ReadFiles| ProcCpu[/proc/cpuinfo]
    DarwinParser -->|CallSysctl| Sysctl[Sysctl]
    DarwinParser -->|CallVMStat| VMStat[VMStat]
    LinuxParser -->|Render| SideBySide[SideBySideRenderer]
    DarwinParser -->|Render| SideBySide
    SideBySide -->|WriteOutput| Stdout[TerminalStdout]
```

## Changelog

Detailed release history is documented in [CHANGELOG.md](CHANGELOG.md).

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

