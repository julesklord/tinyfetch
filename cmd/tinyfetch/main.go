package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func runCommand(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getTerminalWidth() int {
	out := runCommand("tput", "cols")
	if out != "" {
		if w, err := strconv.Atoi(out); err == nil {
			return w
		}
	}
	return 80
}

func runeWidth(r rune) int {
	// Zero-width space, joiners, control chars, variation selectors
	if r == '\u200d' || r == '\u200c' || (r >= '\ufe00' && r <= '\ufe0f') {
		return 0
	}
	// Combining diacritical marks
	if r >= 0x0300 && r <= 0x036F {
		return 0
	}
	// Wide ranges (2 columns)
	// Emojis / Pictographs in SMP (Plane 1): U+1F000 to U+1FAFF
	if r >= 0x1F000 && r <= 0x1FAFF {
		return 2
	}
	// Miscellaneous Symbols and Pictographs, Emoticons, Ornamental Dingbats, etc. in BMP
	if r >= 0x2600 && r <= 0x27BF {
		return 2
	}
	// CJK ranges
	if (r >= 0x2E80 && r <= 0x2FDF) || // CJK Radicals
		(r >= 0x3000 && r <= 0x9FFF) || // Hiragana, Katakana, CJK Unified Ideographs
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility
		(r >= 0xFF01 && r <= 0xFF60) || // Fullwidth Forms
		(r >= 0xFFE0 && r <= 0xFFE6) {
		return 2
	}
	// Default
	return 1
}

func visualLength(s string) int {
	raw := stripANSI(s)
	length := 0
	for _, r := range raw {
		length += runeWidth(r)
	}
	return length
}

func truncateANSI(s string, limit int) string {
	raw := stripANSI(s)
	if visualLength(raw) <= limit {
		return s
	}

	var builder strings.Builder
	visualLen := 0
	inEscape := false
	restoreCode := "\033[0m"
	targetLen := limit - 1
	if targetLen < 0 {
		targetLen = 0
	}

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			builder.WriteByte(s[i])
			continue
		}
		if inEscape {
			builder.WriteByte(s[i])
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}

		if visualLen < targetLen {
			r, size := utf8.DecodeRuneInString(s[i:])
			w := runeWidth(r)
			if visualLen+w <= targetLen {
				builder.WriteRune(r)
				visualLen += w
			} else {
				visualLen = targetLen
			}
			i += size - 1
		}
	}
	builder.WriteString("…")
	builder.WriteString(restoreCode)
	return builder.String()
}

func printJSON(host, osName, kernel, uptime, shell, cpu, mem, disk string, keys, vals []string) {
	fmt.Printf("{\n")
	fmt.Printf("  \"host\": %q,\n", host)
	fmt.Printf("  \"os\": %q,\n", osName)
	fmt.Printf("  \"kernel\": %q,\n", kernel)
	fmt.Printf("  \"uptime\": %q,\n", uptime)
	fmt.Printf("  \"shell\": %q,\n", shell)
	fmt.Printf("  \"cpu\": %q,\n", cpu)
	fmt.Printf("  \"memory\": %q,\n", mem)
	fmt.Printf("  \"disk\": %q", disk)

	if len(keys) > 0 {
		fmt.Printf(",\n  \"plugins\": {\n")
		for i := 0; i < len(keys); i++ {
			cleanVal := stripANSI(vals[i])
			fmt.Printf("    %q: %q", keys[i], cleanVal)
			if i < len(keys)-1 {
				fmt.Printf(",\n")
			} else {
				fmt.Printf("\n")
			}
		}
		fmt.Printf("  }\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("}\n")
}

func printXML(host, osName, kernel, uptime, shell, cpu, mem, disk string, keys, vals []string) {
	fmt.Printf("<tinyfetch>\n")
	fmt.Printf("  <host>%s</host>\n", host)
	fmt.Printf("  <os>%s</os>\n", osName)
	fmt.Printf("  <kernel>%s</kernel>\n", kernel)
	fmt.Printf("  <uptime>%s</uptime>\n", uptime)
	fmt.Printf("  <shell>%s</shell>\n", shell)
	fmt.Printf("  <cpu>%s</cpu>\n", cpu)
	fmt.Printf("  <memory>%s</memory>\n", mem)
	fmt.Printf("  <disk>%s</disk>\n", disk)
	if len(keys) > 0 {
		fmt.Printf("  <plugins>\n")
		for i := 0; i < len(keys); i++ {
			tag := strings.ToLower(keys[i])
			var sb strings.Builder
			for _, r := range tag {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
					sb.WriteRune(r)
				} else {
					sb.WriteRune('_')
				}
			}
			tagStr := sb.String()
			cleanVal := stripANSI(vals[i])
			fmt.Printf("    <%s>%s</%s>\n", tagStr, cleanVal, tagStr)
		}
		fmt.Printf("  </plugins>\n")
	}
	fmt.Printf("</tinyfetch>\n")
}

func printTXT(host, osName, kernel, uptime, shell, cpu, mem, disk string, keys, vals []string) {
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("OS: %s\n", osName)
	fmt.Printf("Kernel: %s\n", kernel)
	fmt.Printf("Uptime: %s\n", uptime)
	fmt.Printf("Shell: %s\n", shell)
	fmt.Printf("CPU: %s\n", cpu)
	fmt.Printf("Memory: %s\n", mem)
	fmt.Printf("Disk: %s\n", disk)
	for i := 0; i < len(keys); i++ {
		fmt.Printf("%s: %s\n", keys[i], stripANSI(vals[i]))
	}
}

func getOSName() string {
	if runtime.GOOS == "darwin" {
		name := runCommand("sw_vers", "-productName")
		ver := runCommand("sw_vers", "-productVersion")
		if name != "" && ver != "" {
			return name + " " + ver
		}
		return "macOS"
	}
	if runtime.GOOS == "linux" {
		file, err := os.Open("/etc/os-release")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					val := strings.TrimPrefix(line, "PRETTY_NAME=")
					return strings.Trim(val, "\"")
				}
			}
		}
		return "Linux"
	}
	return runtime.GOOS
}

func getDistroID() string {
	if runtime.GOOS == "darwin" {
		return "darwin"
	}
	if runtime.GOOS == "linux" {
		file, err := os.Open("/etc/os-release")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "ID=") {
					val := strings.TrimPrefix(line, "ID=")
					return strings.Trim(val, "\"")
				}
			}
		}
	}
	return "linux"
}

func getUptime() string {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/uptime")
		if err == nil {
			parts := strings.Fields(string(data))
			if len(parts) > 0 {
				if sec, err := strconv.ParseFloat(parts[0], 64); err == nil {
					h := int(sec) / 3600
					m := (int(sec) % 3600) / 60
					return fmt.Sprintf("%dh %dm", h, m)
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		out := runCommand("sysctl", "-n", "kern.boottime")
		if out != "" {
			idx := strings.Index(out, "sec = ")
			if idx != -1 {
				s := out[idx+6:]
				comma := strings.Index(s, ",")
				if comma != -1 {
					secStr := strings.TrimSpace(s[:comma])
					if sec, err := strconv.ParseInt(secStr, 10, 64); err == nil {
						diff := time.Now().Unix() - sec
						h := diff / 3600
						m := (diff % 3600) / 60
						return fmt.Sprintf("%dh %dm", h, m)
					}
				}
			}
		}
	}
	// Generic fallback
	uptimeStr := runCommand("uptime")
	if uptimeStr != "" {
		// Very simple parser for uptime output: look for "up"
		idx := strings.Index(uptimeStr, "up ")
		if idx != -1 {
			s := uptimeStr[idx+3:]
			comma := strings.Index(s, ",")
			if comma != -1 {
				return strings.TrimSpace(s[:comma])
			}
		}
	}
	return "n/a"
}

func getCPU() string {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/cpuinfo")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "model name") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		brand := runCommand("sysctl", "-n", "machdep.cpu.brand_string")
		if brand != "" {
			return brand
		}
		model := runCommand("sysctl", "-n", "hw.model")
		if model != "" {
			return model
		}
	}
	return "Unknown CPU"
}

func getMemory() string {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/meminfo")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var total, avail int64
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "MemTotal:") {
					fmt.Sscanf(line, "MemTotal: %d kB", &total)
				} else if strings.HasPrefix(line, "MemAvailable:") {
					fmt.Sscanf(line, "MemAvailable: %d kB", &avail)
				}
			}
			if total > 0 {
				usedPct := (total - avail) * 100 / total
				return fmt.Sprintf("%d%% (%dMB)", usedPct, total/1024)
			}
		}
	} else if runtime.GOOS == "darwin" {
		totalBytesStr := runCommand("sysctl", "-n", "hw.memsize")
		totalBytes, _ := strconv.ParseInt(totalBytesStr, 10, 64)
		totalMB := totalBytes / 1024 / 1024

		pageSizeStr := runCommand("bash", "-c", "vm_stat | awk '/page size of/ {print $8}' | tr -d '.'")
		pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
		if pageSize == 0 {
			pageSize = 4096
		}

		freePagesStr := runCommand("bash", "-c", "vm_stat | awk '/Pages free:/ {print $3}' | tr -d '.'")
		freePages, _ := strconv.ParseInt(freePagesStr, 10, 64)

		inactivePagesStr := runCommand("bash", "-c", "vm_stat | awk '/Pages inactive:/ {print $3}' | tr -d '.'")
		inactivePages, _ := strconv.ParseInt(inactivePagesStr, 10, 64)

		if freePages > 0 && totalMB > 0 {
			freeMB := (freePages + inactivePages) * pageSize / 1024 / 1024
			usedMB := totalMB - freeMB
			pct := usedMB * 100 / totalMB
			return fmt.Sprintf("%d%% (%dMB)", pct, totalMB)
		}
		if totalMB > 0 {
			return fmt.Sprintf("n/a (%dMB)", totalMB)
		}
	}
	return "n/a"
}

func getDisk() string {
	out := runCommand("df", "-h", "/")
	if out != "" {
		lines := strings.Split(out, "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				return fmt.Sprintf("%s (%s)", fields[0], fields[4])
			}
		}
	}
	return "n/a"
}

func getBar(pct int) string {
	filled := pct / 10
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled
	color := "\033[01;32m" // Green
	if pct > 80 {
		color = "\033[01;31m" // Red
	} else if pct > 50 {
		color = "\033[01;33m" // Yellow
	}
	restore := "\033[0m"
	gray := "\033[00;37m"

	var sb strings.Builder
	sb.WriteString(color)
	for i := 0; i < filled; i++ {
		sb.WriteString("█")
	}
	sb.WriteString(restore + gray)
	for i := 0; i < empty; i++ {
		sb.WriteString("░")
	}
	sb.WriteString(restore)
	return sb.String()
}

func stripANSI(s string) string {
	var builder strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		builder.WriteByte(s[i])
	}
	return builder.String()
}

func main() {
	noASCII := false
	minimal := false
	noFrame := false
	outputFmt := ""

	for _, arg := range os.Args[1:] {
		if arg == "--no-ascii" {
			noASCII = true
		} else if arg == "--minimal" {
			minimal = true
		} else if arg == "--noframe" {
			noFrame = true
		} else if strings.HasPrefix(arg, "--output=") {
			outputFmt = strings.TrimPrefix(arg, "--output=")
		} else if arg == "--help" || arg == "-h" {
			fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--output=json|xml|txt]\n", os.Args[0])
			os.Exit(0)
		}
	}

	hostname, _ := os.Hostname()
	osName := getOSName()
	kernel := runCommand("uname", "-r")
	uptimeVal := getUptime()
	shellVal := os.Getenv("SHELL")
	if shellVal == "" {
		shellVal = "sh"
	}
	cpuVal := getCPU()

	// Memory & Progress Bar
	memRaw := getMemory()
	memVal := memRaw
	if strings.Contains(memRaw, "%") {
		pctPart := strings.Split(memRaw, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			memVal = getBar(pct) + " " + memRaw
		}
	}

	// Disk & Progress Bar
	diskRaw := getDisk()
	diskVal := diskRaw
	if strings.Contains(diskRaw, "%") {
		idx := strings.Index(diskRaw, "%")
		start := idx
		for start > 0 && diskRaw[start-1] >= '0' && diskRaw[start-1] <= '9' {
			start--
		}
		if pctStr := diskRaw[start:idx]; pctStr != "" {
			if pct, err := strconv.Atoi(pctStr); err == nil {
				diskVal = getBar(pct) + " " + diskRaw
			}
		}
	}

	// Colors
	restore := "\033[0m"
	lblue := "\033[01;34m"
	lyellow := "\033[01;33m"
	lcyan := "\033[01;36m"
	white := "\033[01;37m"

	// Setup Logo
	var logo []string
	if !noASCII {
		distroID := getDistroID()
		homeDir, _ := os.UserHomeDir()
		// Paths to search
		searchPaths := []string{
			"./ascii/" + distroID + ".txt",
			homeDir + "/.local/share/tinyfetch/ascii/" + distroID + ".txt",
			"/usr/local/share/tinyfetch/ascii/" + distroID + ".txt",
			"/usr/share/tinyfetch/ascii/" + distroID + ".txt",
		}

		asciiPath := ""
		for _, path := range searchPaths {
			if _, err := os.Stat(path); err == nil {
				asciiPath = path
				break
			}
		}

		// Fallback to generic if not found
		if asciiPath == "" {
			fallback := "linux"
			if runtime.GOOS == "darwin" {
				fallback = "darwin"
			}
			fallbackPaths := []string{
				"./ascii/" + fallback + ".txt",
				homeDir + "/.local/share/tinyfetch/ascii/" + fallback + ".txt",
				"/usr/local/share/tinyfetch/ascii/" + fallback + ".txt",
				"/usr/share/tinyfetch/ascii/" + fallback + ".txt",
			}
			for _, path := range fallbackPaths {
				if _, err := os.Stat(path); err == nil {
					asciiPath = path
					break
				}
			}
		}

		if asciiPath != "" {
			file, err := os.Open(asciiPath)
			if err == nil {
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					logo = append(logo, scanner.Text())
				}
			}
		}

		// Hardcoded fallbacks if no file is available
		if len(logo) == 0 {
			if runtime.GOOS == "darwin" {
				logo = []string{
					lcyan + "      .---." + restore,
					lcyan + "     /     \\" + restore,
					lcyan + "     \\__   /" + restore,
					lcyan + "    /   `-' \\" + restore,
					lcyan + "   |         |" + restore,
					lcyan + "    \\       /" + restore,
					lcyan + "     `-...-'" + restore,
				}
			} else {
				logo = []string{
					lyellow + "     .---." + restore,
					lyellow + "    /     \\" + restore,
					lblue + "    \\ " + restore + white + "o o" + restore + lblue + " /" + restore,
					lyellow + "    /  \\-/ \\" + restore,
					lyellow + "   / /     \\ \\" + restore,
					lyellow + "  ( (_     _ ) )" + restore,
					lyellow + "   `(_`---'_)''" + restore,
				}
			}
		}
	}

	// Setup Info
	info := []string{
		lblue + "Host:" + restore + "   " + hostname,
		lblue + "OS:" + restore + "     " + osName,
		lblue + "Kernel:" + restore + " " + kernel,
		lblue + "Uptime:" + restore + " " + uptimeVal,
		lblue + "Shell:" + restore + "  " + shellVal,
		lblue + "CPU:" + restore + "    " + cpuVal,
		lblue + "Memory:" + restore + " " + memVal,
		lblue + "Disk:" + restore + "   " + diskVal,
	}

	var pluginKeys []string
	var pluginVals []string

	// Scan ./plugins directory
	if entries, err := os.ReadDir("./plugins"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				infoPath := "./plugins/" + entry.Name()
				fileInfo, err := entry.Info()
				if err == nil && (fileInfo.Mode()&0111 != 0) {
					out := runCommand(infoPath)
					if out != "" {
						lines := strings.Split(out, "\n")
						pluginOut := strings.TrimSpace(lines[0])
						if pluginOut != "" {
							if strings.Contains(pluginOut, ":") {
								parts := strings.SplitN(pluginOut, ":", 2)
								k := parts[0]
								v := strings.TrimSpace(parts[1])
								pluginKeys = append(pluginKeys, k)
								pluginVals = append(pluginVals, v)
								info = append(info, lblue+k+":"+restore+" "+v)
							} else {
								name := entry.Name()
								if idx := strings.Index(name, "."); idx != -1 {
									name = name[:idx]
								}
								if len(name) > 0 {
									name = strings.ToUpper(name[:1]) + name[1:]
								}
								pluginKeys = append(pluginKeys, name)
								pluginVals = append(pluginVals, pluginOut)
								info = append(info, lblue+name+":"+restore+" "+pluginOut)
							}
						}
					}
				}
			}
		}
	}

	// Intercept output format flag early
	if outputFmt != "" {
		switch outputFmt {
		case "json":
			printJSON(hostname, osName, kernel, uptimeVal, shellVal, cpuVal, memRaw, diskRaw, pluginKeys, pluginVals)
			os.Exit(0)
		case "xml":
			printXML(hostname, osName, kernel, uptimeVal, shellVal, cpuVal, memRaw, diskRaw, pluginKeys, pluginVals)
			os.Exit(0)
		case "txt":
			printTXT(hostname, osName, kernel, uptimeVal, shellVal, cpuVal, memRaw, diskRaw, pluginKeys, pluginVals)
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", outputFmt)
			os.Exit(1)
		}
	}

	// Scan ./plugins/extended directory
	var extInfo []string
	hasExt := false
	if !minimal {
		if entries, err := os.ReadDir("./plugins/extended"); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					infoPath := "./plugins/extended/" + entry.Name()
					fileInfo, err := entry.Info()
					if err == nil && (fileInfo.Mode()&0111 != 0) {
						out, err := exec.Command(infoPath).Output()
						if err == nil {
							rawOut := string(out)
							if strings.TrimSpace(rawOut) != "" {
								lines := strings.Split(rawOut, "\n")
								// Remove trailing empty line caused by Split on final newline
								if len(lines) > 0 && lines[len(lines)-1] == "" {
									lines = lines[:len(lines)-1]
								}
								if len(lines) > 0 {
									for _, line := range lines {
										extInfo = append(extInfo, line)
									}
									extInfo = append(extInfo, "---") // subtle separation token
									hasExt = true
								}
							}
						}
					}
				}
			}
		}
	}
	// Remove trailing separator token if present
	if len(extInfo) > 0 && extInfo[len(extInfo)-1] == "---" {
		extInfo = extInfo[:len(extInfo)-1]
	}

	// Calculate maximum logo raw length
	leftW := 0
	if !noASCII {
		for _, line := range logo {
			rawLen := visualLength(line)
			if rawLen > leftW {
				leftW = rawLen
			}
		}
		if leftW < 16 {
			leftW = 16
		}
	}

	// Calculate maximum info raw length
	rightW := 0
	for _, line := range info {
		rawLen := visualLength(line)
		if rawLen > rightW {
			rightW = rawLen
		}
	}

	// Calculate maximum extended info raw length
	extW := 0
	if hasExt {
		for _, line := range extInfo {
			rawLen := visualLength(line)
			if rawLen > extW {
				extW = rawLen
			}
		}
		if extW < 24 {
			extW = 24
		}
	}

	// Get terminal width
	termW := getTerminalWidth()

	// DYNAMIC RESPONSIVE LAYOUT
	minLogoW := leftW
	if noASCII {
		minLogoW = 0
	}

	// Only disable features if the terminal is physically too small to fit the shrunken columns
	if !noASCII && hasExt {
		if termW < 65 {
			noASCII = true
			minLogoW = 0
		}
	}

	if hasExt {
		if termW < 45 {
			hasExt = false
			extInfo = nil
			extW = 0
		}
	}

	if !noASCII && !hasExt {
		if termW < 41 {
			noASCII = true
			minLogoW = 0
		}
	}

	// Proportional scaling to use the entire terminal width
	totalBorders := 9
	if noASCII {
		totalBorders = 5
	}
	if noFrame {
		totalBorders = 6 // spaces instead of borders
	}

	available := termW - minLogoW - totalBorders
	if hasExt {
		rightW = available * 45 / 100
		extW = available - rightW
		if rightW < 20 {
			rightW = 20
		}
		if extW < 20 {
			extW = 20
		}
	} else {
		rightW = available
		if rightW < 20 {
			rightW = 20
		}
	}

	// Re-evaluate maxLines after scaling/disabling
	maxLines := len(info)
	if !noASCII && len(logo) > maxLines {
		maxLines = len(logo)
	}
	if hasExt && len(extInfo) > maxLines {
		maxLines = len(extInfo)
	}

	borderCol := lblue

	if noFrame {
		// Borderless Rendering
		for i := 0; i < maxLines; i++ {
			logoPrint := ""
			if !noASCII && i < len(logo) {
				logoPrint = logo[i]
			}
			lRaw := visualLength(logoPrint)
			lPadCount := leftW - lRaw
			lPadding := ""
			if lPadCount > 0 {
				lPadding = strings.Repeat(" ", lPadCount)
			}

			infoPrint := ""
			if i < len(info) {
				infoPrint = info[i]
			}
			infoPrint = truncateANSI(infoPrint, rightW)
			rRaw := visualLength(infoPrint)
			rPadCount := rightW - rRaw
			rPadding := ""
			if rPadCount > 0 {
				rPadding = strings.Repeat(" ", rPadCount)
			}

			ePrint := ""
			if hasExt && i < len(extInfo) {
				ePrint = extInfo[i]
			}
			if ePrint == "---" {
				ePrint = "\033[00;37m" + strings.Repeat("╌", extW) + restore
			} else {
				ePrint = truncateANSI(ePrint, extW)
			}

			var sb strings.Builder
			if !noASCII {
				sb.WriteString(" " + logoPrint + lPadding + "   ")
			}
			sb.WriteString(infoPrint + rPadding)
			if hasExt {
				sb.WriteString("   " + ePrint)
			}
			fmt.Println(sb.String())
		}
	} else {
		// Framed Card Rendering
		if !hasExt {
			if noASCII {
				// Case 1: Single pane (Info)
				topLine := borderCol + "┌" + strings.Repeat("─", rightW+2) + "┐" + restore
				botLine := borderCol + "└" + strings.Repeat("─", rightW+2) + "┘" + restore
				fmt.Println(topLine)
				for i := 0; i < maxLines; i++ {
					rLine := ""
					if i < len(info) {
						rLine = info[i]
					}
					rLine = truncateANSI(rLine, rightW)
					rRaw := visualLength(rLine)
					rPadCount := rightW - rRaw
					rPadding := ""
					if rPadCount > 0 {
						rPadding = strings.Repeat(" ", rPadCount)
					}
					fmt.Printf("%s│%s %s%s %s│\n", borderCol, restore, rLine, rPadding, borderCol)
				}
				fmt.Println(botLine)
			} else {
				// Case 2: Double pane (Logo + Info)
				topLine := borderCol + "┌" + strings.Repeat("─", leftW+2) + "┬" + strings.Repeat("─", rightW+2) + "┐" + restore
				botLine := borderCol + "└" + strings.Repeat("─", leftW+2) + "┴" + strings.Repeat("─", rightW+2) + "┘" + restore
				fmt.Println(topLine)
				for i := 0; i < maxLines; i++ {
					logoPrint := ""
					if i < len(logo) {
						logoPrint = logo[i]
					}
					lRaw := visualLength(logoPrint)
					lPadCount := leftW - lRaw
					lPadding := ""
					if lPadCount > 0 {
						lPadding = strings.Repeat(" ", lPadCount)
					}

					infoPrint := ""
					if i < len(info) {
						infoPrint = info[i]
					}
					infoPrint = truncateANSI(infoPrint, rightW)
					rRaw := visualLength(infoPrint)
					rPadCount := rightW - rRaw
					rPadding := ""
					if rPadCount > 0 {
						rPadding = strings.Repeat(" ", rPadCount)
					}

					fmt.Printf("%s│%s %s%s %s│%s %s%s %s│\n",
						borderCol, restore, logoPrint, lPadding,
						borderCol, restore, infoPrint, rPadding,
						borderCol)
				}
				fmt.Println(botLine)
			}
		} else {
			if noASCII {
				// Case 3: Double pane (Info + Extended)
				topLine := borderCol + "┌" + strings.Repeat("─", rightW+2) + "┬" + strings.Repeat("─", extW+2) + "┐" + restore
				botLine := borderCol + "└" + strings.Repeat("─", rightW+2) + "┴" + strings.Repeat("─", extW+2) + "┘" + restore
				fmt.Println(topLine)
				for i := 0; i < maxLines; i++ {
					rLine := ""
					if i < len(info) {
						rLine = info[i]
					}
					rLine = truncateANSI(rLine, rightW)
					rRaw := visualLength(rLine)
					rPadCount := rightW - rRaw
					rPadding := ""
					if rPadCount > 0 {
						rPadding = strings.Repeat(" ", rPadCount)
					}

					eLine := ""
					if i < len(extInfo) {
						eLine = extInfo[i]
					}
					if eLine == "---" {
						eLine = "\033[00;37m" + strings.Repeat("╌", extW) + restore
					} else {
						eLine = truncateANSI(eLine, extW)
					}
					eRaw := visualLength(eLine)
					ePadCount := extW - eRaw
					ePadding := ""
					if ePadCount > 0 {
						ePadding = strings.Repeat(" ", ePadCount)
					}

					fmt.Printf("%s│%s %s%s %s│%s %s%s %s│\n",
						borderCol, restore, rLine, rPadding,
						borderCol, restore, eLine, ePadding,
						borderCol)
				}
				fmt.Println(botLine)
			} else {
				// Case 4: Triple pane (Logo + Info + Extended)
				topLine := borderCol + "┌" + strings.Repeat("─", leftW+2) + "┬" + strings.Repeat("─", rightW+2) + "┬" + strings.Repeat("─", extW+2) + "┐" + restore
				botLine := borderCol + "└" + strings.Repeat("─", leftW+2) + "┴" + strings.Repeat("─", rightW+2) + "┴" + strings.Repeat("─", extW+2) + "┘" + restore
				fmt.Println(topLine)
				for i := 0; i < maxLines; i++ {
					logoPrint := ""
					if i < len(logo) {
						logoPrint = logo[i]
					}
					lRaw := visualLength(logoPrint)
					lPadCount := leftW - lRaw
					lPadding := ""
					if lPadCount > 0 {
						lPadding = strings.Repeat(" ", lPadCount)
					}

					infoPrint := ""
					if i < len(info) {
						infoPrint = info[i]
					}
					infoPrint = truncateANSI(infoPrint, rightW)
					rRaw := visualLength(infoPrint)
					rPadCount := rightW - rRaw
					rPadding := ""
					if rPadCount > 0 {
						rPadding = strings.Repeat(" ", rPadCount)
					}

					ePrint := ""
					if i < len(extInfo) {
						ePrint = extInfo[i]
					}
					if ePrint == "---" {
						ePrint = "\033[00;37m" + strings.Repeat("╌", extW) + restore
					} else {
						ePrint = truncateANSI(ePrint, extW)
					}
					eRaw := visualLength(ePrint)
					ePadCount := extW - eRaw
					ePadding := ""
					if ePadCount > 0 {
						ePadding = strings.Repeat(" ", ePadCount)
					}

					fmt.Printf("%s│%s %s%s %s│%s %s%s %s│%s %s%s %s│\n",
						borderCol, restore, logoPrint, lPadding,
						borderCol, restore, infoPrint, rPadding,
						borderCol, restore, ePrint, ePadding,
						borderCol)
				}
				fmt.Println(botLine)
			}
		}
	}
}
