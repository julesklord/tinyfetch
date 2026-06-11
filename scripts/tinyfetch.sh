#!/usr/bin/env bash
# Robust and portable tinyfetch script
set -euo pipefail

if [ "${1-}" = "--help" ] || [ "${1-}" = "-h" ]; then
  echo "Usage: $0 [--no-ascii] [--minimal] [--noframe] [--output=json|xml|txt]"
  exit 0
fi

NO_ASCII=0
MINIMAL=0
NO_FRAME=0
OUTPUT=""

for a in "$@"; do
  case "$a" in
    --no-ascii) NO_ASCII=1 ;;
    --minimal)  MINIMAL=1 ;;
    --noframe)  NO_FRAME=1 ;;
    --output=*) OUTPUT="${a#*=}" ;;
  esac
done

# Structured output formatters
print_json() {
  local host_esc
  host_esc=$(printf "%s" "$HOST" | sed 's/"/\\"/g')
  local os_esc
  os_esc=$(printf "%s" "$OS_NAME" | sed 's/"/\\"/g')
  local kernel_esc
  kernel_esc=$(printf "%s" "$KERNEL" | sed 's/"/\\"/g')
  local uptime_esc
  uptime_esc=$(printf "%s" "$UPTIME" | sed 's/"/\\"/g')
  local shell_esc
  shell_esc=$(printf "%s" "$SHELL_VAL" | sed 's/"/\\"/g')
  local cpu_esc
  cpu_esc=$(printf "%s" "$CPU" | sed 's/"/\\"/g')
  local mem_esc
  mem_esc=$(printf "%s" "$MEM_RAW" | sed 's/"/\\"/g')
  local disk_esc
  disk_esc=$(printf "%s" "$DISK_RAW" | sed 's/"/\\"/g')

  printf "{\n"
  printf '  "host": "%s",\n' "$host_esc"
  printf '  "os": "%s",\n' "$os_esc"
  printf '  "kernel": "%s",\n' "$kernel_esc"
  printf '  "uptime": "%s",\n' "$uptime_esc"
  printf '  "shell": "%s",\n' "$shell_esc"
  printf '  "cpu": "%s",\n' "$cpu_esc"
  printf '  "memory": "%s",\n' "$mem_esc"
  printf '  "disk": "%s"' "$disk_esc"

  if [ ${#plugin_keys[@]} -gt 0 ]; then
    printf ",\n"
    printf '  "plugins": {\n'
    for ((i=0; i<${#plugin_keys[@]}; i++)); do
      local k="${plugin_keys[i]}"
      local v="${plugin_vals[i]}"
      local v_esc
      v_esc=$(printf "%s" "$v" | sed 's/"/\\"/g' | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
      local k_esc
      k_esc=$(printf "%s" "$k" | sed 's/"/\\"/g')
      printf '    "%s": "%s"' "$k_esc" "$v_esc"
      if [ $i -lt $((${#plugin_keys[@]} - 1)) ]; then
        printf ",\n"
      else
        printf "\n"
      fi
    done
    printf '  }\n'
  else
    printf "\n"
  fi
  printf "}\n"
}

print_xml() {
  printf "<tinyfetch>\n"
  printf "  <host>%s</host>\n" "$HOST"
  printf "  <os>%s</os>\n" "$OS_NAME"
  printf "  <kernel>%s</kernel>\n" "$KERNEL"
  printf "  <uptime>%s</uptime>\n" "$UPTIME"
  printf "  <shell>%s</shell>\n" "$SHELL_VAL"
  printf "  <cpu>%s</cpu>\n" "$CPU"
  printf "  <memory>%s</memory>\n" "$MEM_RAW"
  printf "  <disk>%s</disk>\n" "$DISK_RAW"
  if [ ${#plugin_keys[@]} -gt 0 ]; then
    printf "  <plugins>\n"
    for ((i=0; i<${#plugin_keys[@]}; i++)); do
      local k="${plugin_keys[i]}"
      local v="${plugin_vals[i]}"
      local v_clean
      v_clean=$(printf "%s" "$v" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
      local tag
      tag=$(echo "$k" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/_/g')
      printf "    <%s>%s</%s>\n" "$tag" "$v_clean" "$tag"
    done
    printf "  </plugins>\n"
  fi
  printf "</tinyfetch>\n"
}

print_txt() {
  printf "Host: %s\n" "$HOST"
  printf "OS: %s\n" "$OS_NAME"
  printf "Kernel: %s\n" "$KERNEL"
  printf "Uptime: %s\n" "$UPTIME"
  printf "Shell: %s\n" "$SHELL_VAL"
  printf "CPU: %s\n" "$CPU"
  printf "Memory: %s\n" "$MEM_RAW"
  printf "Disk: %s\n" "$DISK_RAW"
  for ((i=0; i<${#plugin_keys[@]}; i++)); do
    local k="${plugin_keys[i]}"
    local v="${plugin_vals[i]}"
    local v_clean
    v_clean=$(printf "%s" "$v" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
    printf "%s: %s\n" "$k" "$v_clean"
  done
}

# Detection of OS
OS_TYPE=$(uname -s)

# Helper functions for portable resource gathering
get_os_name() {
  if [ "$OS_TYPE" = "Darwin" ]; then
    if command -v sw_vers >/dev/null 2>&1; then
      printf "%s %s" "$(sw_vers -productName)" "$(sw_vers -productVersion)"
    else
      echo "macOS"
    fi
  else
    if [ -f /etc/os-release ]; then
      grep '^PRETTY_NAME' /etc/os-release | cut -d= -f2 | tr -d '"'
    else
      echo "$OS_TYPE"
    fi
  fi
}

get_uptime() {
  if [ -f /proc/uptime ]; then
    local uptime_sec
    uptime_sec=$(cut -d. -f1 /proc/uptime)
    local h=$((uptime_sec / 3600))
    local m=$(((uptime_sec % 3600) / 60))
    echo "${h}h ${m}m"
  elif [ "$OS_TYPE" = "Darwin" ] && command -v sysctl >/dev/null 2>&1; then
    local boot_time
    boot_time=$(sysctl -n kern.boottime 2>/dev/null | awk -F'[=,]' '{print $2}' | tr -d ' ')
    if [ -n "$boot_time" ]; then
      local now
      now=$(date +%s)
      local diff=$((now - boot_time))
      local h=$((diff / 3600))
      local m=$(((diff % 3600) / 60))
      echo "${h}h ${m}m"
    else
      uptime | sed -E 's/^.*up[[:space:]]+([^,]+),.*$/\1/'
    fi
  else
    uptime | sed -E 's/^.*up[[:space:]]+([^,]+),.*$/\1/'
  fi
}

get_cpu() {
  if [ -f /proc/cpuinfo ]; then
    awk -F: '/model name/{print $2; exit}' /proc/cpuinfo | sed 's/^\s*//'
  elif [ "$OS_TYPE" = "Darwin" ] && command -v sysctl >/dev/null 2>&1; then
    sysctl -n machdep.cpu.brand_string 2>/dev/null || sysctl -n hw.model 2>/dev/null || echo "Unknown CPU"
  else
    echo "Unknown CPU"
  fi
}

get_memory() {
  if [ -f /proc/meminfo ]; then
    awk '/MemTotal/ {t=$2} /MemAvailable/ {a=$2} END {if (t>0) printf "%d%% (%dMB)", (t-a)/t*100, t/1024; else print "n/a"}' /proc/meminfo
  elif [ "$OS_TYPE" = "Darwin" ] && command -v sysctl >/dev/null 2>&1 && command -v vm_stat >/dev/null 2>&1; then
    local total_bytes
    total_bytes=$(sysctl -n hw.memsize 2>/dev/null)
    local total_mb=$((total_bytes / 1024 / 1024))
    
    local page_size
    page_size=$(vm_stat | awk '/page size of/ {print $8}' | tr -d '.')
    [ -z "$page_size" ] && page_size=4096
    
    local free_pages
    free_pages=$(vm_stat | awk '/Pages free:/ {print $3}' | tr -d '.')
    local inactive_pages
    inactive_pages=$(vm_stat | awk '/Pages inactive:/ {print $3}' | tr -d '.')
    
    if [ -n "$free_pages" ] && [ -n "$inactive_pages" ]; then
      local free_mb=$(((free_pages + inactive_pages) * page_size / 1024 / 1024))
      local used_mb=$((total_mb - free_mb))
      local pct=$((used_mb * 100 / total_mb))
      echo "${pct}% (${total_mb}MB)"
    else
      echo "n/a (${total_mb}MB)"
    fi
  else
    echo "n/a"
  fi
}

# Colors
ESC=$(printf '\033')
RESTORE="${ESC}[0m"
LBLUE="${ESC}[01;34m"
LYELLOW="${ESC}[01;33m"
LCYAN="${ESC}[01;36m"
WHITE="${ESC}[01;37m"
LRED="${ESC}[01;31m"
LGREEN="${ESC}[01;32m"
LIGHTGRAY="${ESC}[00;37m"

# Progress Bar Helper
get_bar() {
  local pct=$1
  local filled=$((pct / 10))
  [ $filled -gt 10 ] && filled=10
  local empty=$((10 - filled))
  local bar=""
  
  local color="$LGREEN"
  if [ "$pct" -gt 80 ]; then
    color="$LRED"
  elif [ "$pct" -gt 50 ]; then
    color="$LYELLOW"
  fi
  
  bar="${color}"
  for ((i=0; i<filled; i++)); do bar="${bar}█"; done
  bar="${bar}${RESTORE}${LIGHTGRAY}"
  for ((i=0; i<empty; i++)); do bar="${bar}░"; done
  bar="${bar}${RESTORE}"
  echo "$bar"
}

# Resolve values safely
HOST=$(hostname)
OS_NAME=$(get_os_name)
KERNEL=$(uname -r)
UPTIME=$(get_uptime)
SHELL_VAL="${SHELL-sh}"
CPU=$(get_cpu)

# Memory with visual bar
MEM_RAW=$(get_memory)
if [[ "$MEM_RAW" == *"%"* ]]; then
  MEM_PCT=$(echo "$MEM_RAW" | cut -d% -f1)
  MEM_BAR=$(get_bar "$MEM_PCT")
  MEMORY="${MEM_BAR} ${MEM_RAW}"
else
  MEMORY="$MEM_RAW"
fi

# Disk with visual bar
DISK_RAW=$(df -h / | awk 'NR==2 {print $1 " (" $5 ")"}')
DISK_PCT=$(echo "$DISK_RAW" | grep -o '[0-9]\+%' | tr -d '%' || echo "0")
DISK_BAR=$(get_bar "$DISK_PCT")
DISK="${DISK_BAR} ${DISK_RAW}"

# Get Distro ID
get_distro_id() {
  if [ "$OS_TYPE" = "Darwin" ]; then
    echo "darwin"
  elif [ -f /etc/os-release ]; then
    local id
    id=$(grep '^ID=' /etc/os-release | cut -d= -f2 | tr -d '"')
    echo "${id:-linux}"
  else
    echo "linux"
  fi
}

DISTRO_ID=$(get_distro_id)

# Find ASCII file path
ASCII_FILE=""
for path in "./ascii/${DISTRO_ID}.txt" "/usr/local/share/tinyfetch/ascii/${DISTRO_ID}.txt" "/usr/share/tinyfetch/ascii/${DISTRO_ID}.txt"; do
  if [ -f "$path" ]; then
    ASCII_FILE="$path"
    break
  fi
done

# If not found, try fallback linux.txt or darwin.txt
if [ -z "$ASCII_FILE" ]; then
  fallback_name="linux"
  if [ "$OS_TYPE" = "Darwin" ]; then
    fallback_name="darwin"
  fi
  for path in "./ascii/${fallback_name}.txt" "/usr/local/share/tinyfetch/ascii/${fallback_name}.txt" "/usr/share/tinyfetch/ascii/${fallback_name}.txt"; do
    if [ -f "$path" ]; then
      ASCII_FILE="$path"
      break
    fi
  done
fi

logo=()
if [ "$NO_ASCII" -eq 0 ]; then
  if [ -n "$ASCII_FILE" ]; then
    while IFS= read -r line || [ -n "$line" ]; do
      logo+=("$line")
    done < "$ASCII_FILE"
  else
    if [ "$OS_TYPE" = "Darwin" ]; then
      logo[0]="${LCYAN}      .---.${RESTORE}"
      logo[1]="${LCYAN}     /     \\${RESTORE}"
      logo[2]="${LCYAN}     \\__   /${RESTORE}"
      logo[3]="${LCYAN}    /   \`-' \\${RESTORE}"
      logo[4]="${LCYAN}   |         |${RESTORE}"
      logo[5]="${LCYAN}    \\       /${RESTORE}"
      logo[6]="${LCYAN}     \`-...-'${RESTORE}"
      logo[7]=""
    else
      logo[0]="${LYELLOW}     .---.${RESTORE}"
      logo[1]="${LYELLOW}    /     \\${RESTORE}"
      logo[2]="${LBLUE}    \\ ${RESTORE}${WHITE}o o${RESTORE}${LBLUE} /${RESTORE}"
      logo[3]="${LYELLOW}    /  \\-/ \\${RESTORE}"
      logo[4]="${LYELLOW}   / /     \\ \\${RESTORE}"
      logo[5]="${LYELLOW}  ( (_     _ ) )${RESTORE}"
      logo[6]="${LYELLOW}   \`(_\`---'_)''${RESTORE}"
      logo[7]=""
    fi
  fi
fi

info=()
info[0]="${LBLUE}Host:${RESTORE}   $HOST"
info[1]="${LBLUE}OS:${RESTORE}     $OS_NAME"
info[2]="${LBLUE}Kernel:${RESTORE} $KERNEL"
info[3]="${LBLUE}Uptime:${RESTORE} $UPTIME"
info[4]="${LBLUE}Shell:${RESTORE}  $SHELL_VAL"
info[5]="${LBLUE}CPU:${RESTORE}    $CPU"
info[6]="${LBLUE}Memory:${RESTORE} $MEMORY"
info[7]="${LBLUE}Disk:${RESTORE}   $DISK"

plugin_keys=()
plugin_vals=()

# Scan ./plugins directory
if [ -d "./plugins" ]; then
  # Enable nullglob to avoid executing literally `./plugins/*` if folder is empty
  shopt -s nullglob
  for p in ./plugins/*; do
    if [ -x "$p" ] && [ -f "$p" ]; then
      plugin_out=$("$p" 2>/dev/null | head -n 1)
      if [ -n "$plugin_out" ]; then
        if [[ "$plugin_out" == *":"* ]]; then
          p_key=$(echo "$plugin_out" | cut -d: -f1)
          p_val=$(echo "$plugin_out" | cut -d: -f2- | sed 's/^\s*//')
          plugin_keys+=("$p_key")
          plugin_vals+=("$p_val")
          info+=("${LBLUE}${p_key}:${RESTORE} $p_val")
        else
          label=$(basename "$p" | cut -d. -f1 | sed 's/^[0-9]\+-//')
          label="$(tr '[:lower:]' '[:upper:]' <<< "${label:0:1}")${label:1}"
          plugin_keys+=("$label")
          plugin_vals+=("$plugin_out")
          info+=("${LBLUE}${label}:${RESTORE} $plugin_out")
        fi
      fi
    fi
  done
  shopt -u nullglob
fi

# Intercept output format flag early
if [ -n "$OUTPUT" ]; then
  case "$OUTPUT" in
    json) print_json; exit 0 ;;
    xml)  print_xml;  exit 0 ;;
    txt)  print_txt;  exit 0 ;;
    *) echo "Unknown output format: $OUTPUT" >&2; exit 1 ;;
  esac
fi

# Scan ./plugins/extended directory
ext_info=()
HAS_EXT=0
if [ "$MINIMAL" -eq 0 ] && [ -d "./plugins/extended" ]; then
  shopt -s nullglob
  for p in ./plugins/extended/*; do
    if [ -x "$p" ] && [ -f "$p" ]; then
      has_content=0
      # Use temporary array to hold this plugin's output
      tmp_out=()
      while IFS= read -r line || [ -n "$line" ]; do
        tmp_out+=("$line")
        has_content=1
      done < <("$p" 2>/dev/null)
      
      if [ "$has_content" -eq 1 ]; then
        for line in "${tmp_out[@]}"; do
          ext_info+=("$line")
        done
        ext_info+=("---") # subtle separation token
        HAS_EXT=1
      fi
    fi
  done
  shopt -u nullglob
  # Remove trailing separation token if present
  if [ ${#ext_info[@]} -gt 0 ]; then
    last_idx=$((${#ext_info[@]} - 1))
    if [ "${ext_info[$last_idx]}" = "---" ]; then
      unset "ext_info[$last_idx]"
      # Re-index array after unset to avoid sparse index issues
      ext_info=("${ext_info[@]}")
    fi
  fi
fi

# Get terminal width
term_w=$(tput cols 2>/dev/null || echo 80)

# Print with novel card layout
max_lines=${#info[@]}
if [ "$NO_ASCII" -eq 0 ] && [ ${#logo[@]} -gt "$max_lines" ]; then
  max_lines=${#logo[@]}
fi
if [ "$HAS_EXT" -eq 1 ] && [ ${#ext_info[@]} -gt "$max_lines" ]; then
  max_lines=${#ext_info[@]}
fi

# Calculate maximum logo raw length
left_w=0
if [ "$NO_ASCII" -eq 0 ]; then
  for line in "${logo[@]}"; do
    raw=$(printf "%s" "$line" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
    if [ ${#raw} -gt $left_w ]; then
      left_w=${#raw}
    fi
  done
  [ $left_w -lt 16 ] && left_w=16
fi

# Calculate maximum info raw length
right_w=0
for line in "${info[@]}"; do
  raw=$(printf "%s" "$line" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
  if [ ${#raw} -gt $right_w ]; then
    right_w=${#raw}
  fi
done

# Calculate maximum extended info raw length
ext_w=0
if [ "$HAS_EXT" -eq 1 ]; then
  for line in "${ext_info[@]}"; do
    raw=$(printf "%s" "$line" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
    if [ ${#raw} -gt $ext_w ]; then
      ext_w=${#raw}
    fi
  done
  [ $ext_w -lt 24 ] && ext_w=24
fi

# DYNAMIC RESPONSIVE LAYOUT
min_logo_w=$left_w
[ "$NO_ASCII" -eq 1 ] && min_logo_w=0

# Only disable features if the terminal is physically too small to fit the shrunken columns
if [ "$NO_ASCII" -eq 0 ] && [ "$HAS_EXT" -eq 1 ]; then
  if [ "$term_w" -lt 65 ]; then
    NO_ASCII=1
    min_logo_w=0
  fi
fi

if [ "$HAS_EXT" -eq 1 ]; then
  if [ "$term_w" -lt 45 ]; then
    HAS_EXT=0
    ext_info=()
    ext_w=0
  fi
fi

if [ "$NO_ASCII" -eq 0 ] && [ "$HAS_EXT" -eq 0 ]; then
  if [ "$term_w" -lt 41 ]; then
    NO_ASCII=1
    min_logo_w=0
  fi
fi

# Limit maximum pane widths to avoid layout explosions
max_right_w=45
max_ext_w=50
[ $right_w -gt $max_right_w ] && right_w=$max_right_w
if [ "$HAS_EXT" -eq 1 ]; then
  [ $ext_w -gt $max_ext_w ] && ext_w=$max_ext_w
fi

# Proportional scaling down to fit remaining space
total_borders=9
[ "$NO_ASCII" -eq 1 ] && total_borders=5
[ "$NO_FRAME" -eq 1 ] && total_borders=6 # spaces instead of borders

available=$((term_w - min_logo_w - total_borders))
if [ $((right_w + ext_w)) -gt $available ]; then
  if [ "$HAS_EXT" -eq 1 ]; then
    right_w=$((available * 45 / 100))
    ext_w=$((available - right_w))
    [ $right_w -lt 20 ] && right_w=20
    [ $ext_w -lt 20 ] && ext_w=20
  else
    right_w=$available
    [ $right_w -lt 20 ] && right_w=20
  fi
fi

strip_ansi() {
  printf "%s" "$1" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g'
}

truncate_ansi() {
  local str="$1"
  local limit="$2"
  local stripped
  stripped=$(printf "%s" "$str" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
  if [ ${#stripped} -le "$limit" ]; then
    echo "$str"
  else
    echo -e "${stripped:0:$((limit - 1))}…${RESTORE}"
  fi
}

repeat_char() {
  local char="$1"
  local count="$2"
  local out=""
  for ((k=0; k<count; k++)); do
    out="${out}${char}"
  done
  echo -n "$out"
}

BORDER_COLOR="$LBLUE"

# RENDER FLOW
if [ "$NO_FRAME" -eq 1 ]; then
  # Borderless Rendering
  for ((i=0; i<max_lines; i++)); do
    l_line="${logo[i]:-}"
    r_line="${info[i]:-}"
    e_line="${ext_info[i]:-}"

    # Setup Left Logo
    l_raw=$(strip_ansi "$l_line")
    l_pad=$((left_w - ${#l_raw}))
    l_padding=""
    [ $l_pad -gt 0 ] && l_padding=$(printf "%${l_pad}s" "")

    # Setup Middle Info
    r_line=$(truncate_ansi "$r_line" "$right_w")
    r_raw=$(strip_ansi "$r_line")
    r_pad=$((right_w - ${#r_raw}))
    r_padding=""
    [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")

    # Setup Right Extended Info
    if [ "$e_line" = "---" ]; then
      e_line="${LIGHTGRAY}$(repeat_char "╌" $ext_w)${RESTORE}"
    else
      e_line=$(truncate_ansi "$e_line" "$ext_w")
    fi

    # Output line
    out_line=""
    if [ "$NO_ASCII" -eq 0 ]; then
      out_line=" ${l_line}${l_padding}   "
    fi
    out_line="${out_line}${r_line}${r_padding}"
    if [ "$HAS_EXT" -eq 1 ]; then
      out_line="${out_line}   ${e_line}"
    fi
    echo -e "$out_line"
  done
else
  # Framed Card Rendering
  if [ "$HAS_EXT" -eq 0 ]; then
    if [ "$NO_ASCII" -eq 1 ]; then
      # Case 1: Single pane (Info)
      top_line="${BORDER_COLOR}┌$(repeat_char "─" $((right_w + 2)))┐${RESTORE}"
      bot_line="${BORDER_COLOR}└$(repeat_char "─" $((right_w + 2)))┘${RESTORE}"
      echo -e "$top_line"
      for ((i=0; i<max_lines; i++)); do
        r_line="${info[i]:-}"
        r_line=$(truncate_ansi "$r_line" "$right_w")
        r_raw=$(strip_ansi "$r_line")
        r_pad=$((right_w - ${#r_raw}))
        r_padding=""
        [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
        echo -e "${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│"
      done
      echo -e "$bot_line"
    else
      # Case 2: Double pane (Logo + Info)
      top_line="${BORDER_COLOR}┌$(repeat_char "─" $((left_w + 2)))┬$(repeat_char "─" $((right_w + 2)))┐${RESTORE}"
      bot_line="${BORDER_COLOR}└$(repeat_char "─" $((left_w + 2)))┴$(repeat_char "─" $((right_w + 2)))┘${RESTORE}"
      echo -e "$top_line"
      for ((i=0; i<max_lines; i++)); do
        l_line="${logo[i]:-}"
        r_line="${info[i]:-}"
        l_raw=$(strip_ansi "$l_line")
        l_pad=$((left_w - ${#l_raw}))
        l_padding=""
        [ $l_pad -gt 0 ] && l_padding=$(printf "%${l_pad}s" "")
        
        r_line=$(truncate_ansi "$r_line" "$right_w")
        r_raw=$(strip_ansi "$r_line")
        r_pad=$((right_w - ${#r_raw}))
        r_padding=""
        [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
        echo -e "${BORDER_COLOR}│${RESTORE} ${l_line}${l_padding} ${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│"
      done
      echo -e "$bot_line"
    fi
  else
    if [ "$NO_ASCII" -eq 1 ]; then
      # Case 3: Double pane (Info + Extended)
      top_line="${BORDER_COLOR}┌$(repeat_char "─" $((right_w + 2)))┬$(repeat_char "─" $((ext_w + 2)))┐${RESTORE}"
      bot_line="${BORDER_COLOR}└$(repeat_char "─" $((right_w + 2)))┴$(repeat_char "─" $((ext_w + 2)))┘${RESTORE}"
      echo -e "$top_line"
      for ((i=0; i<max_lines; i++)); do
        r_line="${info[i]:-}"
        e_line="${ext_info[i]:-}"
        
        r_line=$(truncate_ansi "$r_line" "$right_w")
        r_raw=$(strip_ansi "$r_line")
        r_pad=$((right_w - ${#r_raw}))
        r_padding=""
        [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
        
        if [ "$e_line" = "---" ]; then
          e_line="${LIGHTGRAY}$(repeat_char "╌" $ext_w)${RESTORE}"
        else
          e_line=$(truncate_ansi "$e_line" "$ext_w")
        fi
        e_raw=$(strip_ansi "$e_line")
        e_pad=$((ext_w - ${#e_raw}))
        e_padding=""
        [ $e_pad -gt 0 ] && e_padding=$(printf "%${e_pad}s" "")
        echo -e "${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│${RESTORE} ${e_line}${e_padding} ${BORDER_COLOR}│"
      done
      echo -e "$bot_line"
    else
      # Case 4: Triple pane (Logo + Info + Extended)
      top_line="${BORDER_COLOR}┌$(repeat_char "─" $((left_w + 2)))┬$(repeat_char "─" $((right_w + 2)))┬$(repeat_char "─" $((ext_w + 2)))┐${RESTORE}"
      bot_line="${BORDER_COLOR}└$(repeat_char "─" $((left_w + 2)))┴$(repeat_char "─" $((right_w + 2)))┴$(repeat_char "─" $((ext_w + 2)))┘${RESTORE}"
      echo -e "$top_line"
      for ((i=0; i<max_lines; i++)); do
        l_line="${logo[i]:-}"
        r_line="${info[i]:-}"
        e_line="${ext_info[i]:-}"
        l_raw=$(strip_ansi "$l_line")
        l_pad=$((left_w - ${#l_raw}))
        l_padding=""
        [ $l_pad -gt 0 ] && l_padding=$(printf "%${l_pad}s" "")
        
        r_line=$(truncate_ansi "$r_line" "$right_w")
        r_raw=$(strip_ansi "$r_line")
        r_pad=$((right_w - ${#r_raw}))
        r_padding=""
        [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
        
        if [ "$e_line" = "---" ]; then
          e_line="${LIGHTGRAY}$(repeat_char "╌" $ext_w)${RESTORE}"
        else
          e_line=$(truncate_ansi "$e_line" "$ext_w")
        fi
        e_raw=$(strip_ansi "$e_line")
        e_pad=$((ext_w - ${#e_raw}))
        e_padding=""
        [ $e_pad -gt 0 ] && e_padding=$(printf "%${e_pad}s" "")
        echo -e "${BORDER_COLOR}│${RESTORE} ${l_line}${l_padding} ${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│${RESTORE} ${e_line}${e_padding} ${BORDER_COLOR}│"
      done
      echo -e "$bot_line"
    fi
  fi
fi

exit 0

