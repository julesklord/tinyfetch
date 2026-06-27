import re

with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

gather_info_func = """
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

"""

# Insert gatherInfo function before main function
content = re.sub(r'func main\(\) \{', gather_info_func + 'func main() {', content)

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
