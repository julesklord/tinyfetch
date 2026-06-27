#!/usr/bin/env bash
# Weather plugin for tinyfetch (queries wttr.in with 1s timeout)
set -euo pipefail

# Try to fetch weather in format: icon + temp (e.g. "☀️ +20°C")
# We use curl with a 1 second timeout to ensure it doesn't block tinyfetch
if ! weather=$(curl -s --connect-timeout 1 "https://wttr.in/?format=%c%t" 2>/dev/null); then
  exit 0
fi

# Clean up output
weather=$(echo "$weather" | xargs)

if [ -n "$weather" ] && [[ "$weather" != *"Error"* ]] && [[ "$weather" != *"Unknown"* ]]; then
  echo "Weather: $weather"
fi
