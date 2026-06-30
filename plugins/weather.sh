#!/usr/bin/env bash
# Weather plugin for arbol (queries wttr.in with 2s timeout)
set -euo pipefail

# Require explicit opt-in for location-based weather fetching
if [ "${ARBOL_ENABLE_WEATHER:-0}" != "1" ]; then
  exit 0
fi

# Try to fetch weather with details: temp/emoji, location, condition
# We use curl with a 2 second timeout to ensure it doesn't block arbol
if ! weather_out=$(curl -s --connect-timeout 2 "https://wttr.in/?format=%c%t\nLocation:+%l\nCondition:+%C" 2>/dev/null); then
  exit 0
fi

# Clean up and print output if valid
if [ -n "$weather_out" ] && [[ "$weather_out" != *"Error"* ]] && [[ "$weather_out" != *"Unknown"* ]]; then
  first_line=$(echo "$weather_out" | head -n 1 | xargs)
  echo "Weather: $first_line"
  echo "$weather_out" | tail -n +2
  
  # Extract temperature number (e.g. +40 or -5)
  temp_num=$(echo "$first_line" | grep -oE '[-+]?[0-9]+' | head -n 1 || echo "")
  if [ -n "$temp_num" ]; then
    # Temperature range: -15°C to 45°C (60 degrees total)
    val=$((temp_num))
    [ "$val" -lt -15 ] && val=-15
    [ "$val" -gt 45 ] && val=45
    # Calculate position (0 to 12)
    pos=$(( (val + 15) * 12 / 60 ))
    bar=""
    for ((i=0; i<=12; i++)); do
      if [ "$i" -eq "$pos" ]; then
        bar="${bar}█"
      else
        bar="${bar}━"
      fi
    done
    echo "Temp Scale: ❄️ ${bar} 🔥"
  fi
fi
