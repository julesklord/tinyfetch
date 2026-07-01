#!/usr/bin/env bash
# Battery status plugin for arbol (cross-platform Linux & macOS)
set -euo pipefail

# ANSI colors
ESC=$(printf '\033')
GREEN="${ESC}[01;32m"
YELLOW="${ESC}[01;33m"
RED="${ESC}[01;31m"
BLUE="${ESC}[01;34m"
RESTORE="${ESC}[0m"

capacity=""
status=""
tech=""
health=""
os_type=$(uname -s)

if [ "$os_type" = "Darwin" ]; then
  if command -v pmset >/dev/null 2>&1; then
    batt_out=$(pmset -g batt 2>/dev/null)
    if echo "$batt_out" | grep -q "InternalBattery"; then
      capacity=$(echo "$batt_out" | grep -o '[0-9]\+%' | tr -d '%')
      if echo "$batt_out" | grep -q "discharging"; then
        status="Discharging"
      else
        status="Charging"
      fi
      # macOS details
      health=$(system_profiler SPPowerDataType 2>/dev/null | grep -i "Condition:" | awk '{print $2}' || echo "")
    fi
  fi
else
  # Linux lookup
  for bat in /sys/class/power_supply/BAT*; do
    if [ -d "$bat" ]; then
      capacity=$(cat "$bat/capacity" 2>/dev/null || echo "")
      status=$(cat "$bat/status" 2>/dev/null || echo "")
      tech=$(cat "$bat/technology" 2>/dev/null || echo "")
      health=$(cat "$bat/status" 2>/dev/null || echo "")
      break
    fi
  done
fi

# Exit if no battery found
if [ -z "$capacity" ]; then
  exit 0
fi

# Format output
icon="🔋"
color="$GREEN"

if [ "$capacity" -le 20 ]; then
  color="$RED"
elif [ "$capacity" -le 50 ]; then
  color="$YELLOW"
fi

if [ "$status" = "Charging" ]; then
  icon="🔌"
  color="$BLUE"
fi

# Draw progress bar
filled=$((capacity / 10))
empty=$((10 - filled))
bar=""
for ((i=0; i<filled; i++)); do bar="${bar}█"; done
for ((i=0; i<empty; i++)); do bar="${bar}░"; done

# Output structure
echo "Battery: ${color}${icon} ${capacity}%${RESTORE} (${status})"
echo "Status: ${status}"
echo "Capacity: ${capacity}%"
echo "Health: ${color}[${bar}]${RESTORE} (${health:-Good})"
if [ -n "$tech" ]; then
  echo "Technology: ${tech}"
fi
