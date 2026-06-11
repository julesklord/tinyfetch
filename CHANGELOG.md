# Changelog

All notable changes to this project will be documented in this file.

## 0.3.1 - Responsive Layout & Structured Output Options
- Implemented a dynamic terminal-responsive resizing system that scales columns, truncates overflow safely, and falls back automatically from 3-pane to 2-pane or 1-pane on smaller terminals to avoid frame deformation.
- Added subtle divider lines (`╌╌╌╌`) between extended plugins.
- Added `--minimal` flag to skip extended plugins entirely.
- Added `--noframe` flag to render borderless columns with clean spacing.
- Added `--output=(json|xml|txt)` to serialize system stats and simple plugins in JSON, XML, or plain TXT format.

## 0.3.0 - Multi-pane Layout & Extended Plugins
- Designed and implemented a flexible layout system that dynamically scales from 1-pane to a 2-pane or 3-pane dashboard based on active assets and plugins.
- Created a new directory `plugins/extended/` to store multi-line, complex plugins.
- Built 3 out-of-the-box extended plugins:
  - Weather Forecast (`plugins/extended/weather_forecast.sh`): Multi-line weather forecast query from wttr.in.
  - Git Commit Graph (`plugins/extended/git_graph.sh`): Displays a beautiful branch history visualization (local or via GitHub API if a token is present).
  - System Dashboard (`plugins/extended/sys_dashboard.sh`): Displays load averages and top memory-consuming processes.
- Enhanced both Shell (`scripts/tinyfetch.sh`) and Go (`cmd/tinyfetch/main.go`) codebases to align, size, and wrap multi-pane boundaries symmetrically using rune counts.

## 0.2.4 - Developer & System plugins
- Created `plugins/battery.sh` to report battery level and charging status across Linux and macOS.
- Created `plugins/ip.sh` to fetch and show both internal IP and public IP addresses (with a 1s connection timeout).
- Created `plugins/k8s.sh` to retrieve active Kubernetes context and namespace using kubectl.

## 0.2.3 - Packages manager plugin
- Created `plugins/packages.sh` to count installed package manager details, supporting native `pacman` packages and foreign (AUR) packages via helper detection (`paru`/`yay`).

## 0.2.2 - Documentation improvements
- Updated `docs/wiki/architecture.md` with system overview, component diagram, and Architecture Decision Records (ADRs) for visual alignment and modular extensibility.
- Updated `docs/wiki/development.md` with a comprehensive guide to writing custom plugins, including constraints, stdout format guidelines, and examples.


## 0.2.1 - Visual alignment fixes & Developer plugins
- Fixed box-drawing layout misalignment in Go by using `utf8.RuneCountInString` instead of byte counts for progress bars and unicode characters.
- Upgraded the Git status plugin with Nerd Fonts, staged/modified/untracked counters, and upstream sync indicators.
- Created 3 new developer plugins: Weather (`plugins/weather.sh`), Docker (`plugins/docker.sh`), and Media Player (`plugins/media.sh`).


## 0.2.0 - Multiplatform stability & Go version
- Refactored `scripts/tinyfetch.sh` for Linux & macOS portability under `set -e`.
- Fixed ShellCheck `SC2034` warning.
- Added support for a modular plugins folder (`./plugins/`).
- Added dynamic distro ASCII logo loading from `ascii/` text files with automatic fallbacks.
- Added visual progress bars for memory and disk usage metrics.
- Replaced standard printing with an innovative double-pane terminal card layout with box-drawing borders.
- Implemented compiled Go version in `cmd/tinyfetch/main.go`.
- Added test suite in `tests/test.sh`.
- Created standard FMG files (`docs/AGENT.md`, `docs/GEMINI.md`, etc.).





## 0.1.0 - Initial skeleton
- Repository scaffolded with standard structure and README.

