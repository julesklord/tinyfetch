#!/usr/bin/env bash
# System Dashboard plugin for tinyfetch (extended pane)
set -euo pipefail

ESC=$(printf '\033')
BLUE="${ESC}[01;34m"
RESTORE="${ESC}[0m"
YELLOW="${ESC}[01;33m"

echo "${BLUE}┌── System Dashboard ──┐${RESTORE}"

# Load average
load=$(uptime | awk -F'load average:' '{print $2}' 2>/dev/null || echo "")
if [ -z "$load" ]; then
  load=$(uptime | awk -F'load averages:' '{print $2}' 2>/dev/null || echo "")
fi
if [ -n "$load" ]; then
  echo "  Load Avg:${YELLOW}${load}${RESTORE}"
fi

# Top Processes by memory
echo "  Top Memory Processes:"
os_type=$(uname -s)
if [ "$os_type" = "Darwin" ]; then
  ps -eo %mem,comm -m 2>/dev/null | head -n 4 | tail -n 3 | while read -r mem comm; do
    if [ -n "$mem" ] && [ -n "$comm" ]; then
      name=$(basename "$comm")
      printf "    %-15s %s%%\n" "$name" "$mem"
    fi
  done
else
  ps -eo %mem,comm --sort=-%mem 2>/dev/null | head -n 4 | tail -n 3 | while read -r mem comm; do
    if [ -n "$mem" ] && [ -n "$comm" ]; then
      name=$(basename "$comm")
      printf "    %-15s %s%%\n" "$name" "$mem"
    fi
  done
fi
