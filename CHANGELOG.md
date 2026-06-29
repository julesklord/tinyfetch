# Changelog

All notable changes to this project will be documented in this file.

## 0.5.0 - Project Renaming, Rich Graphical Plugins & CPU Resources
- Renamed the entire project from `tinyfetch` to `arbol` across modules, file paths, variables, configuration files, and documentation.
- Integrated a real-time CPU usage monitor with a progress bar and CPU temperature detection in default resources.
- Added a horizontal thermometer scale graph (`Temp Scale`) representing weather relative temperature to `weather.sh`.
- Added dynamic, multi-format package manager detection (dpkg, rpm, pacman, brew, flatpak, snap) and a visual distribution ratio bar (`Ratio`) to `packages.sh`.
- Added network traffic counter (Download Rx / Upload Tx in GB/MB/KB) using active interface detection to `ip.sh`.
- Added list of top 3 active containers and uptime stats to `docker.sh`.

## 0.4.0 - Concurrent Plugins, Code Modularization & Security Enhancements
- Refactored monolithic `main()` function in Go implementation into modular helper functions (`parseFlags()`, `gatherInfo()`, and `renderOutput()`).
- Consolidated multiple parameters of export printing functions under a unified `SystemInfo` struct.
- Implemented concurrent execution of simple plugins and extended plugins using `sync.WaitGroup` to avoid sequential blocking.
- Centralized padding calculation into a new helper function `padString()`.
- Secured plugin loading by looking up executable parent directory or `ARBOL_PLUGINS_DIR` environment variable instead of using insecure relative paths.
- Switched wttr.in weather API and icanhazip.com IP API requests to HTTPS to prevent MITM attacks.
- Added unit test coverage for `sysinfo.go` (including timeout/exit codes/sanity checks), `export.go` (with all special XML characters), `stripANSI()`, and `getBar()`.
- Updated GitHub Actions CI workflow to use pinned actions and run shellcheck over the entire repository.
- Removed obsolete scripting helper files and resolved minor shellcheck warnings in shell plugins.

## 0.3.8 - Robust ANSI Truncation, CJK Length & XML Sanitization
- Implemented robust, color-preserving `truncate_ansi` in Bash script to prevent border color loss.
- Added full support for double-width CJK characters, emojis, and zero-width markers in Bash `visual_len`.
- Added strict XML tag sanitization and escaping for both Go and Bash versions.
- Optimized execution paths to bypass visual layout computations when structured formats (JSON/XML/TXT) are requested.
- Implemented process group cleanup in Bash to prevent background processes from leaking on plugin timeouts.
- Allowed help arguments (`-h` / `--help`) at any position in the arguments list in Bash.
- Resolved potential crash conditions under strict mode (`set -e`) by adding fallback operations for missing os-release parameters.
- Clamped progress bar metrics in both versions to prevent drawing overflows.

## 0.3.7 - Code Modularization, Unit Tests & Plugin Timeouts
- Refactored Go implementation in `cmd/arbol/` into multiple single-responsibility files (`main.go`, `render.go`, `sysinfo.go`, `export.go`).
- Created Go table-driven unit test suite in `cmd/arbol/main_test.go` checking string manipulation and unicode width functions.
- Added 2-second timeout limits to simple and extended plugin executions in both Go and Bash.
- Highly optimized visual length calculations and ANSI color stripping in Bash script by replacing external `sed` processes with native string patterns and regular expression loops.

## 0.3.6 - Central Column Scaling & Layout Tweaks
- Increased default central column width proportion to 60% of available terminal width.
- Tuned column width scaling for extended panel 3 to split 50/50 with the central pane.

## 0.3.5 - Terminal Scaling & Paths Expansion
- Expanded columns to fill the full terminal width by default instead of keeping a static maximum layout size.
- Added user local share directory (`~/.local/share/arbol/ascii/`) to logo text search paths.

## 0.3.4 - Weather Alignment Correction
- Fixed box border misalignment when displaying wttr.in responses by utilizing corrected visual length calculations.

## 0.3.3 - Box Border Alignment Fix
- Fixed a layout misalignment bug in the Go implementation of `truncateANSI` where the ellipsis character `…` was appended without subtracting its visual width from the limit. This caused truncated lines to be 1 character wider than expected, pushing the vertical borders outward and deforming the card layout.

## 0.3.2 - Proportional Layout Scaling
- Fixed premature single-pane fallback behavior. The side-by-side logo and info columns are now fully preserved by default using proportional column scaling down to 65, 45, and 41 terminal columns, rather than aggressively dropping logo panels.

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
- Enhanced both Shell (`scripts/arbol.sh`) and Go (`cmd/arbol/main.go`) codebases to align, size, and wrap multi-pane boundaries symmetrically using rune counts.
`
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
- Refactored `scripts/arbol.sh` for Linux & macOS portability under `set -e`.
- Fixed ShellCheck `SC2034` warning.
- Added support for a modular plugins folder (`./plugins/`).
- Added dynamic distro ASCII logo loading from `ascii/` text files with automatic fallbacks.
- Added visual progress bars for memory and disk usage metrics.
- Replaced standard printing with an innovative double-pane terminal card layout with box-drawing borders.
- Implemented compiled Go version in `cmd/arbol/main.go`.
- Added test suite in `tests/test.sh`.
- Created standard FMG files (`docs/AGENT.md`, `docs/GEMINI.md`, etc.).





## 0.1.0 - Initial skeleton
- Repository scaffolded with standard structure and README.
