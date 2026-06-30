#!/usr/bin/env bash
# Git commit graph plugin for arbol (extended pane)
set -euo pipefail

# ANSI colors
ESC=$(printf '\033')
BLUE="${ESC}[01;34m"
RESTORE="${ESC}[0m"



if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  exit 0
fi

echo "${BLUE}┌── Git Commit Graph ──┐${RESTORE}"

# If GITHUB_TOKEN is available, we fetch from GitHub API to show remote commit status
token="${GITHUB_TOKEN:-${GH_TOKEN:-}}"
if [ -n "$token" ]; then
  remote_url=$(git remote get-url origin 2>/dev/null || echo "")
  if [[ "$remote_url" == *"github.com"* ]]; then
    repo=$(echo "$remote_url" | sed -E 's/.*github.com[:\/]([^/]+\/[^.]+).*/\1/')
    response=$(curl -s --connect-timeout 2 -H "Authorization: token $token" "https://api.github.com/repos/$repo/commits?per_page=5" 2>/dev/null || echo "")
    if [ -n "$response" ] && ! echo "$response" | grep -q '"message":'; then
      if command -v jq >/dev/null 2>&1; then
        echo "$response" | jq -r '
          def strip_ansi:
            gsub("\u001b\\[[0-9;?]*[a-zA-Z]"; "") |
            gsub("\u001b\\][^\u0007]*\u0007"; "") |
            gsub("[[:cntrl:]]"; "");
          .[] |
          (.commit.author.name | strip_ansi) as $name |
          (.commit.message | split("\n")[0] | strip_ansi) as $msg |
          "  \u001b[32m●\u001b[0m \(.sha[0:7]) \u001b[33m(\($name))\u001b[0m \($msg)"' | cut -c1-60
        exit 0
      fi
    fi
  fi
fi

# Fallback to local git log graph
git log --graph --oneline --decorate -n 5 --color=always 2>/dev/null | while IFS= read -r line; do
  echo "  $line"
done
