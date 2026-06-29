#!/usr/bin/env bash
# Media status plugin for tinyfetch (queries playerctl if available)
set -euo pipefail

# Check if playerctl exists
if ! command -v playerctl >/dev/null 2>&1; then
  exit 0
fi

# Get player status
status=$(playerctl status 2>/dev/null || echo "")

if [ "$status" = "Playing" ] || [ "$status" = "Paused" ]; then
  artist=$(playerctl metadata artist 2>/dev/null || echo "")
  title=$(playerctl metadata title 2>/dev/null || echo "")
  album=$(playerctl metadata album 2>/dev/null || echo "")
  player=$(playerctl -l 2>/dev/null | head -n 1 || echo "unknown")
  
  # Clean up empty values
  artist=$(echo "$artist" | xargs)
  title=$(echo "$title" | xargs)
  album=$(echo "$album" | xargs)
  player=$(echo "$player" | xargs)
  
  # Build description
  track=""
  if [ -n "$artist" ] && [ -n "$title" ]; then
    track="$artist - $title"
  elif [ -n "$title" ]; then
    track="$title"
  fi
  
  if [ -n "$track" ]; then
    ESC=$(printf '\033')
    GREEN="${ESC}[01;32m"
    YELLOW="${ESC}[01;33m"
    RESTORE="${ESC}[0m"
    
    icon="󰎆"
    color="$YELLOW"
    if [ "$status" = "Playing" ]; then
      icon=""
      color="$GREEN"
    fi
    
    echo "Music: ${color}${icon}${RESTORE} $status"
    echo "Status: $status"
    echo "Track: $track"
    if [ -n "$album" ]; then
      echo "Album: $album"
    fi
    echo "Player: $player"
  fi
fi
