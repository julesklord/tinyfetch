#!/usr/bin/env bash
# Weather plugin for tinyfetch (queries wttr.in with 2s timeout)
set -euo pipefail

# Try to fetch weather with details: temp/emoji, location, condition
# We use curl with a 2 second timeout to ensure it doesn't block tinyfetch
if ! weather_out=$(curl -s --connect-timeout 2 "https://wttr.in/?format=%c%t\nLocation:+%l\nCondition:+%C" 2>/dev/null); then
  exit 0
fi

# Clean up and print output if valid
if [ -n "$weather_out" ] && [[ "$weather_out" != *"Error"* ]] && [[ "$weather_out" != *"Unknown"* ]]; then
  # Split output into lines and print
  echo "Weather: $(echo "$weather_out" | head -n 1 | xargs)"
  echo "$weather_out" | tail -n +2
fi
