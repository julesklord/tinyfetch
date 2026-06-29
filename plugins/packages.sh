#!/usr/bin/env bash
# Package manager count plugin for arbol (cross-platform Linux & macOS)
set -euo pipefail

native_mgr=""
native_count=0
aur_count=0
helper="none"
brew_count=0
flatpak_count=0
snap_count=0

# Detect native managers
if command -v pacman >/dev/null 2>&1; then
  native_count=$(pacman -Qn 2>/dev/null | wc -l | xargs)
  native_mgr="pacman"
  
  # AUR
  aur_count=$(pacman -Qm 2>/dev/null | wc -l | xargs)
  if command -v paru >/dev/null 2>&1; then
    helper="paru"
  elif command -v yay >/dev/null 2>&1; then
    helper="yay"
  else
    helper="pacman"
  fi
elif command -v dpkg-query >/dev/null 2>&1; then
  native_count=$(dpkg-query -f '${binary:Package}\n' -W 2>/dev/null | wc -l | xargs)
  native_mgr="dpkg"
elif command -v rpm >/dev/null 2>&1; then
  native_count=$(rpm -qa 2>/dev/null | wc -l | xargs)
  native_mgr="rpm"
fi

# Homebrew
if command -v brew >/dev/null 2>&1; then
  brew_count=$(brew list --formula 2>/dev/null | wc -l | xargs)
fi

# Flatpak
if command -v flatpak >/dev/null 2>&1; then
  flatpak_count=$(flatpak list 2>/dev/null | wc -l | xargs)
  if [ "$flatpak_count" -gt 0 ]; then
    flatpak_count=$((flatpak_count - 1))
  fi
fi

# Snap
if command -v snap >/dev/null 2>&1; then
  snap_count=$(snap list 2>/dev/null | tail -n +2 | wc -l | xargs || echo "0")
fi

total=$((native_count + aur_count + brew_count + flatpak_count + snap_count))

if [ "$total" -eq 0 ]; then
  exit 0
fi

# Visual distribution bar (10 blocks)
native_bar=0
aur_bar=0
other_bar=0
if [ "$total" -gt 0 ]; then
  native_bar=$(( native_count * 10 / total ))
  aur_bar=$(( aur_count * 10 / total ))
  other_bar=$(( (brew_count + flatpak_count + snap_count) * 10 / total ))
fi

bar=""
for ((i=0; i<native_bar; i++)); do bar="${bar}█"; done
for ((i=0; i<aur_bar; i++)); do bar="${bar}▒"; done
for ((i=0; i<other_bar; i++)); do bar="${bar}░"; done
while [ ${#bar} -lt 10 ]; do
  bar="${bar}░"
done

echo "Packages: $total total"
if [ "$native_count" -gt 0 ]; then
  echo "Native: $native_count ($native_mgr)"
fi
if [ "$aur_count" -gt 0 ]; then
  echo "AUR: $aur_count ($helper)"
fi
if [ "$brew_count" -gt 0 ]; then
  echo "Homebrew: $brew_count"
fi
if [ "$flatpak_count" -gt 0 ]; then
  echo "Flatpak: $flatpak_count"
fi
if [ "$snap_count" -gt 0 ]; then
  echo "Snap: $snap_count"
fi

echo "Ratio: [${bar:0:10}] (native/aur/others)"
