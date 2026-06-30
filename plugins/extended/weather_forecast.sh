#!/usr/bin/env bash
# Weather forecast plugin for arbol (extended pane)
set -euo pipefail

# ANSI colors
ESC=$(printf '\033')
BLUE="${ESC}[01;34m"
RESTORE="${ESC}[0m"

# Require explicit opt-in for location-based weather fetching
if [ "${ARBOL_ENABLE_WEATHER:-0}" != "1" ]; then
  exit 0
fi

# Fetch weather forecast from wttr.in with a 2s timeout
weather=$(curl -s --connect-timeout 2 "https://wttr.in/?0&Q&T" 2>/dev/null || echo "")

if [ -z "$weather" ] || echo "$weather" | grep -q -E "Error|502|503|504|Gateway|Timeout"; then
  exit 0
fi

echo "${BLUE}┌── Weather Forecast ──┐${RESTORE}"
echo "$weather" | while IFS= read -r line; do
  # Avoid printing completely empty trailing lines
  if [ -n "$line" ]; then
    echo "  $line"
  fi
done
