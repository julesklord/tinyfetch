package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TreeNode struct {
	Text     string
	Children []*TreeNode
}

func parseFlags() (bool, bool, bool, string, string, string, string, string, bool, int, string, bool, int) {
	noASCII := false
	minimal := false
	noFrame := false
	outputFmt := ""
	logoMode := "banner" // default is banner
	themeName := ""
	barStyleName := ""
	treeStyleName := ""
	sparklineEnabled := false
	sparklineWidth := 20
	sparklineStyleName := ""
	liveEnabled := false
	liveInterval := 1000 // ms

	for _, arg := range os.Args[1:] {
		if arg == "--no-ascii" {
			noASCII = true
		} else if arg == "--minimal" {
			minimal = true
		} else if arg == "--noframe" {
			noFrame = true
		} else if strings.HasPrefix(arg, "--output=") {
			outputFmt = strings.TrimPrefix(arg, "--output=")
		} else if strings.HasPrefix(arg, "--logo=") {
			logoMode = strings.TrimPrefix(arg, "--logo=")
		} else if strings.HasPrefix(arg, "--theme=") {
			themeName = strings.TrimPrefix(arg, "--theme=")
		} else if strings.HasPrefix(arg, "--bar-style=") {
			barStyleName = strings.TrimPrefix(arg, "--bar-style=")
		} else if strings.HasPrefix(arg, "--tree-style=") {
			treeStyleName = strings.TrimPrefix(arg, "--tree-style=")
		} else if strings.HasPrefix(arg, "--sparkline") {
			if arg == "--sparkline" {
				sparklineEnabled = true
			} else if strings.HasPrefix(arg, "--sparkline=") {
				sparklineEnabled = true
				val := strings.TrimPrefix(arg, "--sparkline=")
				if w, err := strconv.Atoi(val); err == nil && w > 0 {
					sparklineWidth = w
				}
			}
		} else if strings.HasPrefix(arg, "--sparkline-style=") {
			sparklineStyleName = strings.TrimPrefix(arg, "--sparkline-style=")
		} else if strings.HasPrefix(arg, "--live") {
			liveEnabled = true
			if arg == "--live" {
				liveInterval = 1000
			} else if strings.HasPrefix(arg, "--live=") {
				val := strings.TrimPrefix(arg, "--live=")
				if ms, err := strconv.Atoi(val); err == nil && ms > 0 {
					liveInterval = ms
				}
			}
		} else if arg == "--help" || arg == "-h" {
			fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--logo=simple|banner] [--output=json|xml|txt] [--theme=NAME] [--bar-style=STYLE] [--tree-style=STYLE] [--sparkline[=WIDTH]] [--sparkline-style=STYLE] [--live[=MS]]\n", os.Args[0])
			fmt.Println("  Themes: default, catppuccin, catppuccin-mocha, catppuccin-latte, dracula, nord, tokyonight, gruvbox, everforest, monokai, rose-pine, solarized")
			fmt.Println("  Bar styles: block, braille, gradient, dot")
			fmt.Println("  Tree styles: default, rounded, heavy, double, ascii, dotted")
			fmt.Println("  Sparkline styles: block, braille, dots")
			fmt.Println("  --sparkline[=WIDTH] enables inline sparklines (default width: 20)")
			fmt.Println("  --live[=MS] enables live updating mode (default interval: 1000ms)")
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--logo=simple|banner] [--output=json|xml|txt]\n", os.Args[0])
			os.Exit(1)
		}
	}
	return noASCII, minimal, noFrame, outputFmt, logoMode, themeName, barStyleName, treeStyleName, sparklineEnabled, sparklineWidth, sparklineStyleName, liveEnabled, liveInterval
}

func parseBarStyle(name string) BarStyle {
	name = strings.ToLower(name)
	switch name {
	case "block":
		return BarStyleBlock
	case "braille":
		return BarStyleBraille
	case "gradient":
		return BarStyleGradient
	case "dot":
		return BarStyleDot
	default:
		return BarStyleBraille
	}
}

func parseTreeStyle(name string) TreeStyle {
	name = strings.ToLower(name)
	switch name {
	case "default":
		return TreeStyleDefault
	case "rounded":
		return TreeStyleRounded
	case "heavy":
		return TreeStyleHeavy
	case "double":
		return TreeStyleDouble
	case "ascii":
		return TreeStyleASCII
	case "dotted":
		return TreeStyleDotted
	default:
		return TreeStyleDefault
	}
}

func parseSparklineStyle(name string) SparklineStyle {
	name = strings.ToLower(name)
	switch name {
	case "block":
		return SparklineBlock
	case "braille":
		return SparklineBraille
	case "dots":
		return SparklineDots
	default:
		return SparklineBlock
	}
}

func gatherInfo(pluginsDir string) SystemInfo {
	hostname, _ := os.Hostname()
	osName := getOSName()
	kernel := runCommand("uname", "-r")
	uptimeVal := getUptime()
	shellVal := os.Getenv("SHELL")
	if shellVal == "" {
		shellVal = "sh"
	}
	cpuVal := getCPU()
	gpuVal := getGPU()
	dewmVal := getDEWM()
	termVal := getTerminal()

	memRaw := getMemory()
	swapRaw := getSwap()
	diskRaw := getDisk()
	procVal := getProcesses()

	var plugins []PluginInfo

	// Scan plugins directory
	if entries, err := os.ReadDir(pluginsDir); err == nil {
		type pluginResult struct {
			key     string
			val     string
			details []string
			ok      bool
		}
		results := make([]pluginResult, len(entries))
		var wg sync.WaitGroup

		for i, entry := range entries {
			if !entry.IsDir() {
				infoPath := filepath.Join(pluginsDir, entry.Name())
				fileInfo, err := entry.Info()
				if err == nil && (fileInfo.Mode()&0111 != 0) {
					wg.Add(1)
					go func(idx int, path string, name string) {
						defer wg.Done()
						out := runCommandWithTimeout(2*time.Second, path)
						if out != "" {
							lines := strings.Split(out, "\n")
							var cleanLines []string
							for _, l := range lines {
								trimmed := strings.TrimSpace(l)
								if trimmed != "" {
									cleanLines = append(cleanLines, trimmed)
								}
							}
							if len(cleanLines) > 0 {
								firstLine := cleanLines[0]
								var k, v string
								if strings.Contains(firstLine, ":") {
									parts := strings.SplitN(firstLine, ":", 2)
									k = parts[0]
									v = strings.TrimSpace(parts[1])
								} else {
									parsedName := name
									if dotIdx := strings.Index(parsedName, "."); dotIdx != -1 {
										parsedName = parsedName[:dotIdx]
									}
									if len(parsedName) > 0 {
										parsedName = strings.ToUpper(parsedName[:1]) + parsedName[1:]
									}
									k = parsedName
									v = firstLine
								}
								results[idx] = pluginResult{
									key:     k,
									val:     v,
									details: cleanLines[1:],
									ok:      true,
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
				plugins = append(plugins, PluginInfo{
					Key:     res.key,
					Val:     res.val,
					Details: res.details,
				})
			}
		}
	}

	cpuUsageVal := getCPUUsage()
	cpuTempVal := getCPUTemp()

	return SystemInfo{
		Host:      hostname,
		OSName:    osName,
		Kernel:    kernel,
		Uptime:    uptimeVal,
		Shell:     shellVal,
		CPU:       cpuVal,
		GPU:       gpuVal,
		DEWM:      dewmVal,
		Terminal:  termVal,
		Memory:    memRaw,
		Swap:      swapRaw,
		Disk:      diskRaw,
		Processes: procVal,
		CPUUsage:  cpuUsageVal,
		CPUTemp:   cpuTempVal,
		Plugins:   plugins,
	}
}

func formatPluginName(filename string) string {
	name := filename
	if idx := strings.Index(name, "."); idx != -1 {
		name = name[:idx]
	}
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func loadASCIILogo() []string {
	distroID := getDistroID()
	homeDir, _ := os.UserHomeDir()

	exe, err := os.Executable()
	var exeDir string
	if err == nil {
		if realExe, err := filepath.EvalSymlinks(exe); err == nil {
			exeDir = filepath.Dir(realExe)
		} else {
			exeDir = filepath.Dir(exe)
		}
	}

	searchPaths := []string{
		"./ascii/" + distroID + ".txt",
	}
	if exeDir != "" {
		searchPaths = append(searchPaths, filepath.Join(exeDir, "ascii", distroID+".txt"))
	}
	searchPaths = append(searchPaths,
		homeDir+"/.local/share/arbol/ascii/"+distroID+".txt",
		"/usr/local/share/arbol/ascii/"+distroID+".txt",
		"/usr/share/arbol/ascii/"+distroID+".txt",
	)

	asciiPath := ""
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			asciiPath = path
			break
		}
	}

	// Fallback to generic OS file
	if asciiPath == "" {
		fallback := "linux"
		if runtime.GOOS == "darwin" {
			fallback = "darwin"
		}
		fallbackPaths := []string{
			"./ascii/" + fallback + ".txt",
		}
		if exeDir != "" {
			fallbackPaths = append(fallbackPaths, filepath.Join(exeDir, "ascii", fallback+".txt"))
		}
		fallbackPaths = append(fallbackPaths,
			homeDir+"/.local/share/arbol/ascii/"+fallback+".txt",
			"/usr/local/share/arbol/ascii/"+fallback+".txt",
			"/usr/share/arbol/ascii/"+fallback+".txt",
		)
		for _, path := range fallbackPaths {
			if _, err := os.Stat(path); err == nil {
				asciiPath = path
				break
			}
		}
	}

	var logo []string
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

	// Dynamic/hardcoded fallback if file not found
	if len(logo) == 0 {
		if runtime.GOOS == "darwin" {
			logo = []string{
				"\033[96m      .---.\033[0m",
				"\033[96m     /     \\\033[0m",
				"\033[96m     \\__   /\033[0m",
				"\033[96m    /   `-' \\\033[0m",
				"\033[96m   |         |\033[0m",
				"\033[96m    \\       /\033[0m",
				"\033[96m     `-...-'\033[0m",
			}
		} else {
			logo = []string{
				"\033[33m     .---.\033[0m",
				"\033[33m    /     \\\033[0m",
				"\033[34m    \\ \033[0m\033[1;37mo o\033[0m\033[34m /\033[0m",
				"\033[33m    /  \\-/ \\\033[0m",
				"\033[33m   / /     \\ \\\033[0m",
				"\033[33m  ( (_     _ ) )\033[0m",
				"\033[33m   `(_`---'_)''\033[0m",
			}
		}
	}
	return logo
}

func loadASCIILogoBanner() []string {
	distroID := getDistroID()
	homeDir, _ := os.UserHomeDir()

	exe, err := os.Executable()
	var exeDir string
	if err == nil {
		if realExe, err := filepath.EvalSymlinks(exe); err == nil {
			exeDir = filepath.Dir(realExe)
		} else {
			exeDir = filepath.Dir(exe)
		}
	}

	searchPaths := []string{
		"./ascii/" + distroID + "_banner.txt",
	}
	if exeDir != "" {
		searchPaths = append(searchPaths, filepath.Join(exeDir, "ascii", distroID+"_banner.txt"))
	}
	searchPaths = append(searchPaths,
		homeDir+"/.local/share/arbol/ascii/"+distroID+"_banner.txt",
		"/usr/local/share/arbol/ascii/"+distroID+"_banner.txt",
		"/usr/share/arbol/ascii/"+distroID+"_banner.txt",
	)

	asciiPath := ""
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			asciiPath = path
			break
		}
	}

	// Fallback to generic linux_banner/darwin_banner
	if asciiPath == "" {
		fallback := "linux"
		if runtime.GOOS == "darwin" {
			fallback = "darwin"
		}
		fallbackPaths := []string{
			"./ascii/" + fallback + "_banner.txt",
		}
		if exeDir != "" {
			fallbackPaths = append(fallbackPaths, filepath.Join(exeDir, "ascii", fallback+"_banner.txt"))
		}
		fallbackPaths = append(fallbackPaths,
			homeDir+"/.local/share/arbol/ascii/"+fallback+"_banner.txt",
			"/usr/local/share/arbol/ascii/"+fallback+"_banner.txt",
			"/usr/share/arbol/ascii/"+fallback+"_banner.txt",
		)
		for _, path := range fallbackPaths {
			if _, err := os.Stat(path); err == nil {
				asciiPath = path
				break
			}
		}
	}

	var logo []string
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
	return logo
}

func printTree(node *TreeNode, prefixes []string, isLast bool) {
	connectors := GetConnectors()
	theme := GetTheme()
	treeColor := theme.TreeLines

	if len(prefixes) > 0 {
		for _, p := range prefixes[:len(prefixes)-1] {
			fmt.Print(treeColor + p + "\033[0m")
		}
		if isLast {
			fmt.Print(treeColor + connectors.LastBranch + "\033[0m")
		} else {
			fmt.Print(treeColor + connectors.Branch + "\033[0m")
		}
	}
	fmt.Println(node.Text)

	for i, child := range node.Children {
		var nextPrefixes []string
		if len(prefixes) > 0 {
			nextPrefixes = append(nextPrefixes, prefixes...)
			if isLast {
				nextPrefixes[len(nextPrefixes)-1] = connectors.Space
			} else {
				nextPrefixes[len(nextPrefixes)-1] = connectors.Vertical
			}
		}
		nextPrefixes = append(nextPrefixes, connectors.Vertical)
		printTree(child, nextPrefixes, i == len(node.Children)-1)
	}
}

func gradientString(s string, r1, g1, b1, r2, g2, b2 int) string {
	runes := []rune(s)
	n := len(runes)
	if n <= 1 {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r1, g1, b1, s)
	}

	var sb strings.Builder
	for i, r := range runes {
		ratio := float64(i) / float64(n-1)
		currR := int(float64(r1) + ratio*float64(r2-r1))
		currG := int(float64(g1) + ratio*float64(g2-g1))
		currB := int(float64(b1) + ratio*float64(b2-b1))
		sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c", currR, currG, currB, r))
	}
	sb.WriteString("\033[0m")
	return sb.String()
}

func drawBannerLogo(osName string) {
	theme := GetTheme()
	osUpper := strings.ToUpper(osName)
	osUpper = stripANSI(osUpper)
	content := "  A R B O L  //  " + osUpper + "  "
	width := len(content) + 2

	top := "╔" + strings.Repeat("═", width-2) + "╗"
	mid := "║" + content + "║"
	bot := "╚" + strings.Repeat("═", width-2) + "╝"

	start := theme.BannerGradient[0]
	end := theme.BannerGradient[1]

	fmt.Println(gradientString(top, start[0], start[1], start[2], end[0], end[1], end[2]))
	fmt.Println(gradientString(mid, start[0], start[1], start[2], end[0], end[1], end[2]))
	fmt.Println(gradientString(bot, start[0], start[1], start[2], end[0], end[1], end[2]))
}

func renderOutput(noASCII, minimal, noFrame bool, outputFmt string, infoObj SystemInfo, extPluginsDir, logoMode string, sparklineEnabled bool) {
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

	// Swap & Progress Bar
	swapVal := infoObj.Swap
	if strings.Contains(infoObj.Swap, "%") {
		pctPart := strings.Split(infoObj.Swap, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			swapVal = getBar(pct) + " " + infoObj.Swap
		}
	}
	// CPU Usage & Progress Bar
	cpuUsageVal := infoObj.CPUUsage
	if strings.Contains(infoObj.CPUUsage, "%") {
		pctPart := strings.Split(infoObj.CPUUsage, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			cpuUsageVal = getBar(pct) + " " + infoObj.CPUUsage
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

	theme := GetTheme()

	// Styling tokens
	bold := "\033[1m"
	italic := "\033[3m"
	reset := "\033[0m"
	lblue := theme.Secondary
	lcyan := theme.Primary

	// Render Logo/Banner Header at the top
	if !noASCII {
		switch logoMode {
		case "simple":
			logo := loadASCIILogo()
			for _, line := range logo {
				fmt.Println(line)
			}
			if len(logo) > 0 {
				fmt.Println()
			}
		case "banner":
			logo := loadASCIILogoBanner()
			if len(logo) > 0 {
				for _, line := range logo {
					fmt.Println(gradientString(line, theme.BannerGradient[0][0], theme.BannerGradient[0][1], theme.BannerGradient[0][2], theme.BannerGradient[1][0], theme.BannerGradient[1][1], theme.BannerGradient[1][2]))
				}
				fmt.Println()
			} else {
				drawBannerLogo(infoObj.OSName)
				fmt.Println()
			}
		default:
			// Fallback/Default
			logo := loadASCIILogoBanner()
			if len(logo) > 0 {
				for _, line := range logo {
					fmt.Println(gradientString(line, theme.BannerGradient[0][0], theme.BannerGradient[0][1], theme.BannerGradient[0][2], theme.BannerGradient[1][0], theme.BannerGradient[1][1], theme.BannerGradient[1][2]))
				}
				fmt.Println()
			} else {
				drawBannerLogo(infoObj.OSName)
				fmt.Println()
			}
		}
	}

	// Build Tree Root
	titleText := infoObj.Host + " @ " + infoObj.OSName
	rootText := bold + lcyan + "● " + reset + bold + gradientString(titleText, theme.BannerGradient[0][0], theme.BannerGradient[0][1], theme.BannerGradient[0][2], theme.BannerGradient[1][0], theme.BannerGradient[1][1], theme.BannerGradient[1][2])

	root := &TreeNode{
		Text: rootText,
	}

	// Specs category
	specsNode := &TreeNode{Text: lcyan + bold + "⚙ specs" + reset}
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "📦 kernel: " + reset + italic + infoObj.Kernel + reset})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "⏱ uptime: " + reset + italic + infoObj.Uptime + reset})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "💻 shell: " + reset + italic + infoObj.Shell + reset})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "🧠 cpu: " + reset + italic + infoObj.CPU + reset})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "🎮 gpu: " + reset + italic + infoObj.GPU + reset})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "🖥 de/wm: " + reset + italic + infoObj.DEWM + reset})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "📟 terminal: " + reset + italic + infoObj.Terminal + reset})
	root.Children = append(root.Children, specsNode)

	// Resources category
	resourcesNode := &TreeNode{Text: lcyan + bold + "📊 resources" + reset}

	// CPU Usage with optional sparkline
	cpuText := lblue + "📈 cpu usage: " + reset + cpuUsageVal
	if sparklineEnabled {
		cpuText += "  " + getSparklineCPU(theme)
	}
	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: cpuText})

	if infoObj.CPUTemp != "n/a" {
		resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: lblue + "🌡️ cpu temp: " + reset + italic + infoObj.CPUTemp + reset})
	}

	// Memory with optional sparkline
	memText := lblue + "💾 memory: " + reset + memVal
	if sparklineEnabled {
		memText += "  " + getSparklineMem(theme)
	}
	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: memText})

	// Swap with optional sparkline
	swapText := lblue + "🔄 swap: " + reset + swapVal
	if sparklineEnabled {
		swapText += "  " + getSparklineSwap(theme)
	}
	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: swapText})

	// Disk with optional sparkline
	diskText := lblue + "💿 disk: " + reset + diskVal
	if sparklineEnabled {
		diskText += "  " + getSparklineDisk(theme)
	}
	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: diskText})

	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: lblue + "⚡ processes: " + reset + italic + infoObj.Processes + reset})
	root.Children = append(root.Children, resourcesNode)

	// Simple Plugins category
	if len(infoObj.Plugins) > 0 {
		pluginsNode := &TreeNode{Text: lcyan + bold + "🔌 plugins" + reset}
		for _, plug := range infoObj.Plugins {
			key := strings.ToLower(plug.Key)
			val := plug.Val
			plugNode := &TreeNode{Text: lblue + key + ": " + reset + val}
			for _, det := range plug.Details {
				plugNode.Children = append(plugNode.Children, &TreeNode{Text: det})
			}
			pluginsNode.Children = append(pluginsNode.Children, plugNode)
		}
		root.Children = append(root.Children, pluginsNode)
	}

	// Diagnostics category (extended plugins)
	if !minimal {
		if entries, err := os.ReadDir(extPluginsDir); err == nil {
			type extResult struct {
				name  string
				lines []string
				ok    bool
			}
			results := make([]extResult, len(entries))
			var wg sync.WaitGroup

			for i, entry := range entries {
				if !entry.IsDir() {
					infoPath := filepath.Join(extPluginsDir, entry.Name())
					fileInfo, err := entry.Info()
					if err == nil && (fileInfo.Mode()&0111 != 0) {
						wg.Add(1)
						go func(idx int, path string, filename string) {
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
										results[idx] = extResult{
											name:  formatPluginName(filename),
											lines: lines,
											ok:    true,
										}
									}
								}
							}
						}(i, infoPath, entry.Name())
					}
				}
			}
			wg.Wait()

			var diagChildren []*TreeNode
			for _, res := range results {
				if res.ok {
					pluginNode := &TreeNode{Text: lblue + strings.ToLower(res.name) + reset}
					for _, line := range res.lines {
						pluginNode.Children = append(pluginNode.Children, &TreeNode{Text: line})
					}
					diagChildren = append(diagChildren, pluginNode)
				}
			}

			if len(diagChildren) > 0 {
				diagNode := &TreeNode{Text: lcyan + bold + "🔍 diagnostics" + reset}
				diagNode.Children = diagChildren
				root.Children = append(root.Children, diagNode)
			}
		}
	}

	// Render the Tree
	printTree(root, []string{}, true)
}

func runLiveMode(noASCII, minimal, noFrame bool, outputFmt, logoMode string, intervalMs int) {
	interval := time.Duration(intervalMs) * time.Millisecond

	// Initialize sparklines for live mode
	initSparklines(20, interval/2)
	time.Sleep(100 * time.Millisecond)

	pluginsDir := getPluginsDir()
	extPluginsDir := filepath.Join(pluginsDir, "extended")

	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[H\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			infoObj := gatherInfo(pluginsDir)
			// Move cursor to top-left
			fmt.Print("\033[H")
			renderOutput(noASCII, minimal, noFrame, outputFmt, infoObj, extPluginsDir, logoMode, true)
		}
	}
}

func getPluginsDir() string {
	if env := os.Getenv("ARBOL_PLUGINS_DIR"); env != "" {
		return env
	}
	exe, err := os.Executable()
	if err != nil {
		return "./plugins"
	}
	realExe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		realExe = exe
	}
	return filepath.Join(filepath.Dir(realExe), "plugins")
}

func main() {
	noASCII, minimal, noFrame, outputFmt, logoMode, themeName, barStyleName, treeStyleName, sparklineEnabled, sparklineWidth, sparklineStyleName, liveEnabled, liveInterval := parseFlags()

	if themeName != "" {
		if !SetTheme(themeName) {
			fmt.Fprintf(os.Stderr, "Unknown theme: %s\n", themeName)
			os.Exit(1)
		}
	}

	if barStyleName != "" {
		style := parseBarStyle(barStyleName)
		SetBarStyle(style)
	}

	if treeStyleName != "" {
		style := parseTreeStyle(treeStyleName)
		SetTreeStyle(style)
	}

	if sparklineStyleName != "" {
		style := parseSparklineStyle(sparklineStyleName)
		SetSparklineStyle(style)
	}

	if sparklineEnabled {
		initSparklines(sparklineWidth, 500*time.Millisecond)
		// Give sparklines a moment to collect initial data
		time.Sleep(100 * time.Millisecond)
	}

	if liveEnabled {
		runLiveMode(noASCII, minimal, noFrame, outputFmt, logoMode, liveInterval)
		return
	}

	pluginsDir := getPluginsDir()
	extPluginsDir := filepath.Join(pluginsDir, "extended")
	infoObj := gatherInfo(pluginsDir)
	renderOutput(noASCII, minimal, noFrame, outputFmt, infoObj, extPluginsDir, logoMode, sparklineEnabled)
}
