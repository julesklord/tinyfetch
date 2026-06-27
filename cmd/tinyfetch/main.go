package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SystemInfo struct {
	Hostname   string
	OSName     string
	Kernel     string
	UptimeVal  string
	ShellVal   string
	CPUVal     string
	MemRaw     string
	DiskRaw    string
	PluginKeys []string
	PluginVals []string
}

func parseFlags() (bool, bool, bool, string) {
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
	return noASCII, minimal, noFrame, outputFmt
}

func gatherInfo() SystemInfo {
	hostname, _ := os.Hostname()
	osName := getOSName()
	kernel := runCommand("uname", "-r")
	uptimeVal := getUptime()
	shellVal := os.Getenv("SHELL")
	if shellVal == "" {
		shellVal = "sh"
	}
	cpuVal := getCPU()

	memRaw := getMemory()
	diskRaw := getDisk()

	var pluginKeys []string
	var pluginVals []string

	// Scan ./plugins directory
	if entries, err := os.ReadDir("./plugins"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				infoPath := "./plugins/" + entry.Name()
				fileInfo, err := entry.Info()
				if err == nil && (fileInfo.Mode()&0111 != 0) {
					out := runCommandWithTimeout(2*time.Second, infoPath)
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
							}
						}
					}
				}
			}
		}
	}

	return SystemInfo{
		Hostname:   hostname,
		OSName:     osName,
		Kernel:     kernel,
		UptimeVal:  uptimeVal,
		ShellVal:   shellVal,
		CPUVal:     cpuVal,
		MemRaw:     memRaw,
		DiskRaw:    diskRaw,
		PluginKeys: pluginKeys,
		PluginVals: pluginVals,
	}
}

func renderOutput(noASCII, minimal, noFrame bool, outputFmt string, infoObj SystemInfo) {
	// Intercept output format flag early
	if outputFmt != "" {
		switch outputFmt {
		case "json":
			printJSON(infoObj.Hostname, infoObj.OSName, infoObj.Kernel, infoObj.UptimeVal, infoObj.ShellVal, infoObj.CPUVal, infoObj.MemRaw, infoObj.DiskRaw, infoObj.PluginKeys, infoObj.PluginVals)
			os.Exit(0)
		case "xml":
			printXML(infoObj.Hostname, infoObj.OSName, infoObj.Kernel, infoObj.UptimeVal, infoObj.ShellVal, infoObj.CPUVal, infoObj.MemRaw, infoObj.DiskRaw, infoObj.PluginKeys, infoObj.PluginVals)
			os.Exit(0)
		case "txt":
			printTXT(infoObj.Hostname, infoObj.OSName, infoObj.Kernel, infoObj.UptimeVal, infoObj.ShellVal, infoObj.CPUVal, infoObj.MemRaw, infoObj.DiskRaw, infoObj.PluginKeys, infoObj.PluginVals)
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", outputFmt)
			os.Exit(1)
		}
	}

	// Memory & Progress Bar
	memVal := infoObj.MemRaw
	if strings.Contains(infoObj.MemRaw, "%") {
		pctPart := strings.Split(infoObj.MemRaw, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			memVal = getBar(pct) + " " + infoObj.MemRaw
		}
	}

	// Disk & Progress Bar
	diskVal := infoObj.DiskRaw
	if strings.Contains(infoObj.DiskRaw, "%") {
		idx := strings.Index(infoObj.DiskRaw, "%")
		start := idx
		for start > 0 && infoObj.DiskRaw[start-1] >= '0' && infoObj.DiskRaw[start-1] <= '9' {
			start--
		}
		if pctStr := infoObj.DiskRaw[start:idx]; pctStr != "" {
			if pct, err := strconv.Atoi(pctStr); err == nil {
				diskVal = getBar(pct) + " " + infoObj.DiskRaw
			}
		}
	}

	// Colors
	restore := "[0m"
	lblue := "[01;34m"
	lyellow := "[01;33m"
	lcyan := "[01;36m"
	white := "[01;37m"

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
		lblue + "Host:" + restore + "   " + infoObj.Hostname,
		lblue + "OS:" + restore + "     " + infoObj.OSName,
		lblue + "Kernel:" + restore + " " + infoObj.Kernel,
		lblue + "Uptime:" + restore + " " + infoObj.UptimeVal,
		lblue + "Shell:" + restore + "  " + infoObj.ShellVal,
		lblue + "CPU:" + restore + "    " + infoObj.CPUVal,
		lblue + "Memory:" + restore + " " + memVal,
		lblue + "Disk:" + restore + "   " + diskVal,
	}

	for i := 0; i < len(infoObj.PluginKeys); i++ {
		info = append(info, lblue+infoObj.PluginKeys[i]+":"+restore+" "+infoObj.PluginVals[i])
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
						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
						out, err := exec.CommandContext(ctx, infoPath).Output()
						cancel()
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
			// Keep features active, but they will stack vertically if termW is small
		}
	}

	// Calculate the maximum raw lengths for layout decisions
	maxInfoLineW := 0
	for _, line := range info {
		rawLen := visualLength(line)
		if rawLen > maxInfoLineW {
			maxInfoLineW = rawLen
		}
	}

	maxExtLineW := 0
	if hasExt {
		for _, line := range extInfo {
			rawLen := visualLength(line)
			if rawLen > maxExtLineW {
				maxExtLineW = rawLen
			}
		}
	}

	// Determine margins/borders needed for side-by-side layout
	numActivePanes := 0
	if !noASCII {
		numActivePanes++
	}
	numActivePanes++ // info is always active
	if hasExt {
		numActivePanes++
	}

	bordersSideBySide := 5 // 1 pane (noASCII, !hasExt) -> ┌──────┐ (5 chars)
	if numActivePanes == 2 {
		bordersSideBySide = 6 // ┌────┬────┐ (6 chars)
	} else if numActivePanes == 3 {
		bordersSideBySide = 9 // ┌────┬────┬────┐ (9 chars)
	}
	if noFrame {
		bordersSideBySide = 0
		if !noASCII {
			bordersSideBySide += 4 // Logo spacer: 4 spaces
		}
		if hasExt {
			bordersSideBySide += 3 // Extended spacer: 3 spaces
		}
	}

	minInfoW := 20
	minExtW := 0
	if hasExt {
		minExtW = 20
	}
	minSideBySideWidth := leftW + minInfoW + minExtW + bordersSideBySide

	useVerticalLayout := getTerminalWidth() < minSideBySideWidth

	borderCol := lblue

	if useVerticalLayout {
		// Vertical stacked layout: Draw title bars for each pane sequentially
		boxW := termW - 4
		if boxW < 12 {
			boxW = 12
		}

		drawVerticalBox := func(title string, lines []string) {
			if noFrame {
				fmt.Println(borderCol + "--- " + title + " ---" + restore)
				for _, line := range lines {
					printLine := line
					if printLine == "---" {
						printLine = "[00;37m" + strings.Repeat("╌", termW) + restore
					} else {
						printLine = truncateANSI(printLine, termW)
					}
					fmt.Println(printLine)
				}
				fmt.Println()
			} else {
				titleStr := " " + title + " "
				titleLen := len(title) + 2
				fillW := boxW - titleLen - 2
				if fillW < 2 {
					fillW = 2
				}

				topBorder := borderCol + "┌──" + restore + titleStr + borderCol + strings.Repeat("─", fillW) + "┐" + restore
				botBorder := borderCol + "└" + strings.Repeat("─", boxW) + "┘" + restore

				fmt.Println(topBorder)
				for _, line := range lines {
					printLine := line
					if printLine == "---" {
						printLine = "[00;37m" + strings.Repeat("╌", boxW-2) + restore
					} else {
						printLine = truncateANSI(printLine, boxW-2)
					}
					visualLen := visualLength(printLine)
					padding := boxW - 2 - visualLen
					if padding < 0 {
						padding = 0
					}
					padStr := strings.Repeat(" ", padding)
					fmt.Printf("%s│%s %s%s %s│\n", borderCol, restore, printLine, padStr, borderCol)
				}
				fmt.Println(botBorder)
			}
		}

		if !noASCII && len(logo) > 0 {
			drawVerticalBox("OS Logo", logo)
		}
		if len(info) > 0 {
			drawVerticalBox("System Info", info)
		}
		if hasExt && len(extInfo) > 0 {
			drawVerticalBox("Plugins & Diagnostics", extInfo)
		}
	} else {
		// Proportional scaling to use the entire terminal width
		available := termW - minLogoW - bordersSideBySide
		if hasExt {
			rightW = available * 50 / 100
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
					ePrint = "[00;37m" + strings.Repeat("╌", extW) + restore
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
							eLine = "[00;37m" + strings.Repeat("╌", extW) + restore
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
							ePrint = "[00;37m" + strings.Repeat("╌", extW) + restore
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
}

func main() {
	noASCII, minimal, noFrame, outputFmt := parseFlags()
	infoObj := gatherInfo()
	renderOutput(noASCII, minimal, noFrame, outputFmt, infoObj)
}
