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
	"sync"
	"time"
)

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
		type pluginResult struct {
			key string
			val string
			ok  bool
		}
		results := make([]pluginResult, len(entries))
		var wg sync.WaitGroup

		for i, entry := range entries {
			if !entry.IsDir() {
				infoPath := "./plugins/" + entry.Name()
				fileInfo, err := entry.Info()
				if err == nil && (fileInfo.Mode()&0111 != 0) {
					wg.Add(1)
					go func(idx int, path string, name string) {
						defer wg.Done()
						out := runCommandWithTimeout(2*time.Second, path)
						if out != "" {
							lines := strings.Split(out, "\n")
							pluginOut := strings.TrimSpace(lines[0])
							if pluginOut != "" {
								if strings.Contains(pluginOut, ":") {
									parts := strings.SplitN(pluginOut, ":", 2)
									k := parts[0]
									v := strings.TrimSpace(parts[1])
									results[idx] = pluginResult{key: k, val: v, ok: true}
								} else {
									parsedName := name
									if dotIdx := strings.Index(parsedName, "."); dotIdx != -1 {
										parsedName = parsedName[:dotIdx]
									}
									if len(parsedName) > 0 {
										parsedName = strings.ToUpper(parsedName[:1]) + parsedName[1:]
									}
									results[idx] = pluginResult{key: parsedName, val: pluginOut, ok: true}
								}
							}
						}
					}(i, infoPath, entry.Name())
				}
			}
		}
		wg.Wait()
		for _, res := range results {
			if res.ok {
				pluginKeys = append(pluginKeys, res.key)
				pluginVals = append(pluginVals, res.val)
			}
		}
	}

	return SystemInfo{
		Host:   hostname,
		OSName: osName,
		Kernel: kernel,
		Uptime: uptimeVal,
		Shell:  shellVal,
		CPU:    cpuVal,
		Memory: memRaw,
		Disk:   diskRaw,
		Keys:   pluginKeys,
		Vals:   pluginVals,
	}
}

func renderOutput(noASCII, minimal, noFrame bool, outputFmt string, infoObj SystemInfo) {
	// Intercept output format flag early
	if outputFmt != "" {
		switch outputFmt {
		case "json":
			printJSON(infoObj)
			os.Exit(0)
		case "xml":
			printXML(infoObj)
			os.Exit(0)
		case "txt":
			printTXT(infoObj)
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", outputFmt)
			os.Exit(1)
		}
	}

	// Memory & Progress Bar
	memVal := infoObj.Memory
	if strings.Contains(infoObj.Memory, "%") {
		pctPart := strings.Split(infoObj.Memory, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			memVal = getBar(pct) + " " + infoObj.Memory
		}
	}

	// Disk & Progress Bar
	diskVal := infoObj.Disk
	if strings.Contains(infoObj.Disk, "%") {
		idx := strings.Index(infoObj.Disk, "%")
		start := idx
		for start > 0 && infoObj.Disk[start-1] >= '0' && infoObj.Disk[start-1] <= '9' {
			start--
		}
		if pctStr := infoObj.Disk[start:idx]; pctStr != "" {
			if pct, err := strconv.Atoi(pctStr); err == nil {
				diskVal = getBar(pct) + " " + infoObj.Disk
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
		lblue + "Host:" + restore + "   " + infoObj.Host,
		lblue + "OS:" + restore + "     " + infoObj.OSName,
		lblue + "Kernel:" + restore + " " + infoObj.Kernel,
		lblue + "Uptime:" + restore + " " + infoObj.Uptime,
		lblue + "Shell:" + restore + "  " + infoObj.Shell,
		lblue + "CPU:" + restore + "    " + infoObj.CPU,
		lblue + "Memory:" + restore + " " + memVal,
		lblue + "Disk:" + restore + "   " + diskVal,
	}

	for i := 0; i < len(infoObj.Keys); i++ {
		info = append(info, lblue+infoObj.Keys[i]+":"+restore+" "+infoObj.Vals[i])
	}

	// Scan ./plugins/extended directory
	var extInfo []string
	hasExt := false
	if !minimal {
		if entries, err := os.ReadDir("./plugins/extended"); err == nil {
			type extResult struct {
				lines []string
				ok    bool
			}
			results := make([]extResult, len(entries))
			var wg sync.WaitGroup

			for i, entry := range entries {
				if !entry.IsDir() {
					infoPath := "./plugins/extended/" + entry.Name()
					fileInfo, err := entry.Info()
					if err == nil && (fileInfo.Mode()&0111 != 0) {
						wg.Add(1)
						go func(idx int, path string) {
							defer wg.Done()
							ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
							out, err := exec.CommandContext(ctx, path).Output()
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
										results[idx] = extResult{lines: lines, ok: true}
									}
								}
							}
						}(i, infoPath)
					}
				}
			}
			wg.Wait()

			for _, res := range results {
				if res.ok {
					for _, line := range res.lines {
						extInfo = append(extInfo, line)
					}
					extInfo = append(extInfo, "---") // subtle separation token
					hasExt = true
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
						printLine = "\033[00;37m" + strings.Repeat("╌", termW) + restore
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
						printLine = "\033[00;37m" + strings.Repeat("╌", boxW-2) + restore
					} else {
						printLine = truncateANSI(printLine, boxW-2)
					}
					padStr := padString(printLine, boxW-2)
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
				lPadding := padString(logoPrint, leftW)

				infoPrint := ""
				if i < len(info) {
					infoPrint = info[i]
				}
				infoPrint = truncateANSI(infoPrint, rightW)
				rPadding := padString(infoPrint, rightW)

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
						rPadding := padString(rLine, rightW)
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
						lPadding := padString(logoPrint, leftW)

						infoPrint := ""
						if i < len(info) {
							infoPrint = info[i]
						}
						infoPrint = truncateANSI(infoPrint, rightW)
						rPadding := padString(infoPrint, rightW)

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
						rPadding := padString(rLine, rightW)

						eLine := ""
						if i < len(extInfo) {
							eLine = extInfo[i]
						}
						if eLine == "---" {
							eLine = "\033[00;37m" + strings.Repeat("╌", extW) + restore
						} else {
							eLine = truncateANSI(eLine, extW)
						}
						ePadding := padString(eLine, extW)

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
						lPadding := padString(logoPrint, leftW)

						infoPrint := ""
						if i < len(info) {
							infoPrint = info[i]
						}
						infoPrint = truncateANSI(infoPrint, rightW)
						rPadding := padString(infoPrint, rightW)

						ePrint := ""
						if i < len(extInfo) {
							ePrint = extInfo[i]
						}
						if ePrint == "---" {
							ePrint = "\033[00;37m" + strings.Repeat("╌", extW) + restore
						} else {
							ePrint = truncateANSI(ePrint, extW)
						}
						ePadding := padString(ePrint, extW)

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
