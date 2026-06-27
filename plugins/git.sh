#!/usr/bin/env bash
# Advanced Git status plugin for tinyfetch with Nerd Fonts and colors
set -euo pipefail

# Check if inside a git repo
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  exit 0
fi

# Colors
ESC=$(printf '\033')
PURPLE="${ESC}[01;35m"
RESTORE="${ESC}[0m"
GREEN="${ESC}[01;32m"
YELLOW="${ESC}[01;33m"
RED="${ESC}[01;31m"

# Get branch name
branch=$(git branch --show-current 2>/dev/null || git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "detached")

# Get status stats
status_out=$(git status --porcelain 2>/dev/null)

staged=0
modified=0
untracked=0

while IFS= read -r line; do
  [ -z "$line" ] && continue
  index_status="${line:0:1}"
  work_status="${line:1:1}"
  
  if [ "$index_status" = "?" ] && [ "$work_status" = "?" ]; then
    untracked=$((untracked + 1))
  else
    if [ "$index_status" != " " ] && [ "$index_status" != "?" ]; then
      staged=$((staged + 1))
    fi
    if [ "$work_status" != " " ] && [ "$work_status" != "?" ]; then
      modified=$((modified + 1))
    fi
  fi
done <<< "$status_out"

# Get ahead/behind status
ahead=0
behind=0
if git rev-parse --abbrev-ref @\{u\} >/dev/null 2>&1; then
  upstream_status=$(git rev-list --left-right --count HEAD...@\{u\} 2>/dev/null)
  ahead=$(echo "$upstream_status" | cut -f1)
  behind=$(echo "$upstream_status" | cut -f2)
fi

# Build output
output="${PURPLE} ${branch}${RESTORE}"

# Ahead / Behind status
if [ "$ahead" -gt 0 ] || [ "$behind" -gt 0 ]; then
  output="${output} ("
  if [ "$ahead" -gt 0 ]; then
    output="${output}${GREEN}⇡${ahead}${RESTORE}"
  fi
  if [ "$behind" -gt 0 ]; then
    [ "$ahead" -gt 0 ] && output="${output} "
    output="${output}${RED}⇣${behind}${RESTORE}"
  fi
  output="${output})"
fi

# Dirty files status
if [ "$staged" -gt 0 ] || [ "$modified" -gt 0 ] || [ "$untracked" -gt 0 ]; then
  output="${output} ✗"
  details=""
  if [ "$staged" -gt 0 ]; then
    details="${GREEN}●${staged}${RESTORE}"
  fi
  if [ "$modified" -gt 0 ]; then
    [ -n "$details" ] && details="${details} "
    details="${details}${YELLOW}✚${modified}${RESTORE}"
  fi
  if [ "$untracked" -gt 0 ]; then
    [ -n "$details" ] && details="${details} "
    details="${details}${RED}…${untracked}${RESTORE}"
  fi
  output="${output} [${details}]"
else
  output="${output} ${GREEN}✔${RESTORE}"
fi

echo "Git: $output"
