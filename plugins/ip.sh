#!/usr/bin/env bash
# Network IP status plugin for tinyfetch
set -euo pipefail

# ANSI colors
ESC=$(printf '\033')
BLUE="${ESC}[01;34m"
RESTORE="${ESC}[0m"

# Get Local IP (cross-platform fallback)
local_ip=""
if command -v ip >/dev/null 2>&1; then
  local_ip=$(ip route get 1.1.1.1 2>/dev/null | awk '{print $7; exit}' || echo "")
fi

if [ -z "$local_ip" ] && command -v hostname >/dev/null 2>&1; then
  local_ip=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "")
fi

if [ -z "$local_ip" ] && command -v ifconfig >/dev/null 2>&1; then
  local_ip=$(ifconfig 2>/dev/null | grep -E "inet " | grep -v "127.0.0.1" | awk '{print $2; exit}' | tr -d 'addr:' || echo "")
fi

# Get Public IP (with 1s timeout to prevent hanging)
public_ip=$(curl -s --connect-timeout 1 icanhazip.com 2>/dev/null | xargs || echo "")

# Build output
output=""
if [ -n "$local_ip" ]; then
  output="Local: $local_ip"
fi

if [ -n "$public_ip" ]; then
  if [ -n "$output" ]; then
    output="$output | "
  fi
  output="${output}Public: $public_ip"
fi

# Exit if no IP found
if [ -z "$output" ]; then
  exit 0
fi

echo "Network: ${BLUE}󰖩${RESTORE} $output"
