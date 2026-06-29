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
  local_ip=$(ifconfig 2>/dev/null | grep -E "inet " | grep -v "127.0.0.1" | awk '{print $2; exit}' | sed "s/addr://" || echo "")
fi

# Get Public IP (with 1s timeout to prevent hanging)
public_ip=$(curl -s --connect-timeout 1 https://icanhazip.com 2>/dev/null | xargs || echo "")

# Get DNS Server
dns_server=""
if [ -f /etc/resolv.conf ]; then
  dns_server=$(grep -E '^nameserver' /etc/resolv.conf | awk '{print $2}' | head -n 2 | xargs || echo "")
fi

# Get default gateway
gateway=""
if command -v ip >/dev/null 2>&1; then
  gateway=$(ip route show 2>/dev/null | grep default | awk '{print $3}' | head -n 1 || echo "")
fi

# Build output
status="disconnected"
if [ -n "$local_ip" ] || [ -n "$public_ip" ]; then
  status="connected"
fi

echo "Network: ${BLUE}󰖩${RESTORE} $status"
echo "Local IP: ${local_ip:-n/a}"
echo "Public IP: ${public_ip:-n/a}"
if [ -n "$dns_server" ]; then
  echo "DNS: $dns_server"
fi
if [ -n "$gateway" ]; then
  echo "Gateway: $gateway"
fi
