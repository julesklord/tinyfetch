#!/usr/bin/env bash
# Network IP status plugin for arbol
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

# Get default gateway and interface
gateway=""
iface=""
if command -v ip >/dev/null 2>&1; then
  gateway_info=$(ip route show 2>/dev/null | grep default | head -n 1 || echo "")
  gateway=$(echo "$gateway_info" | awk '{print $3}')
  iface=$(echo "$gateway_info" | awk '{print $5}')
else
  # macOS
  gateway=$(route -n get default 2>/dev/null | grep gateway | awk '{print $2}' || echo "")
  iface=$(route -n get default 2>/dev/null | grep interface | awk '{print $2}' || echo "")
fi

# Get Rx/Tx traffic stats
rx_bytes=""
tx_bytes=""
os_type=$(uname -s)

if [ "$os_type" = "Linux" ] && [ -n "$iface" ]; then
  rx_bytes=$(cat "/sys/class/net/$iface/statistics/rx_bytes" 2>/dev/null || echo "")
  tx_bytes=$(cat "/sys/class/net/$iface/statistics/tx_bytes" 2>/dev/null || echo "")
elif [ "$os_type" = "Darwin" ] && [ -n "$iface" ]; then
  stats=$(netstat -I "$iface" -b 2>/dev/null | tail -n 1 || echo "")
  rx_bytes=$(echo "$stats" | awk '{print $7}')
  tx_bytes=$(echo "$stats" | awk '{print $10}')
fi

bytes_to_human() {
  local b=$1
  if [ -z "$b" ] || [ "$b" -eq 0 ]; then
    echo "0 B"
    return
  fi
  if [ "$b" -ge 1073741824 ]; then
    awk "BEGIN {printf \"%.2f GB\n\", $b / 1073741824}"
  elif [ "$b" -ge 1048576 ]; then
    awk "BEGIN {printf \"%.2f MB\n\", $b / 1048576}"
  elif [ "$b" -ge 1024 ]; then
    awk "BEGIN {printf \"%.2f KB\n\", $b / 1024}"
  else
    echo "${b} B"
  fi
}

# Build output
status="disconnected"
if [ -n "$local_ip" ] || [ -n "$public_ip" ]; then
  status="connected"
fi

echo "Network: ${BLUE}󰖩${RESTORE} $status"
echo "Interface: ${iface:-n/a}"
echo "Local IP: ${local_ip:-n/a}"
echo "Public IP: ${public_ip:-n/a}"
if [ -n "$dns_server" ]; then
  echo "DNS: $dns_server"
fi
if [ -n "$gateway" ]; then
  echo "Gateway: $gateway"
fi

if [ -n "$rx_bytes" ] && [ -n "$tx_bytes" ]; then
  rx_human=$(bytes_to_human "$rx_bytes")
  tx_human=$(bytes_to_human "$tx_bytes")
  echo "Download (Rx): 📥 $rx_human"
  echo "Upload (Tx): 📤 $tx_human"
fi
