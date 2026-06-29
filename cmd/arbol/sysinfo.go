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

func runCommand(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func runCommandWithTimeout(timeout time.Duration, name string, arg ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, arg...).Output()
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
	out := runCommand("df", "-Ph", "/")
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

func getGPU() string {
	if runtime.GOOS == "darwin" {
		out := runCommand("bash", "-c", "system_profiler SPDisplaysDataType | grep 'Chipset Model'")
		if out != "" {
			parts := strings.Split(out, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	} else if runtime.GOOS == "linux" {
		out := runCommand("bash", "-c", "lspci | grep -i 'vga\\|3d\\|display'")
		if out != "" {
			lines := strings.Split(out, "\n")
			line := lines[0]
			if idx := strings.Index(line, "controller:"); idx != -1 {
				line = line[idx+11:]
			} else if idx := strings.Index(line, "VGA compatible controller: "); idx != -1 {
				line = line[idx+27:]
			} else if idx := strings.Index(line, "3D controller: "); idx != -1 {
				line = line[idx+15:]
			}
			if idx := strings.Index(line, " (rev "); idx != -1 {
				line = line[:idx]
			}
			return strings.TrimSpace(line)
		}
	}
	return "n/a"
}

func getDEWM() string {
	if runtime.GOOS == "darwin" {
		return "Aqua"
	}
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de != "" {
		return de
	}
	wm := os.Getenv("DESKTOP_SESSION")
	if wm != "" {
		return wm
	}
	return "n/a"
}

func getTerminal() string {
	termProg := os.Getenv("TERM_PROGRAM")
	if termProg != "" {
		return termProg
	}
	termEnv := os.Getenv("TERM")
	if termEnv != "" {
		return termEnv
	}
	return "n/a"
}

func getSwap() string {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/meminfo")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var total, free int64
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "SwapTotal:") {
					fmt.Sscanf(line, "SwapTotal: %d kB", &total)
				} else if strings.HasPrefix(line, "SwapFree:") {
					fmt.Sscanf(line, "SwapFree: %d kB", &free)
				}
			}
			if total > 0 {
				used := total - free
				pct := used * 100 / total
				return fmt.Sprintf("%d%% (%dMB / %dMB)", pct, used/1024, total/1024)
			}
		}
	} else if runtime.GOOS == "darwin" {
		out := runCommand("sysctl", "-n", "vm.swapusage")
		if out != "" {
			return out
		}
	}
	return "n/a"
}

func getProcesses() string {
	if runtime.GOOS == "linux" {
		files, err := os.ReadDir("/proc")
		if err == nil {
			count := 0
			for _, f := range files {
				if f.IsDir() {
					if _, err := strconv.Atoi(f.Name()); err == nil {
						count++
					}
				}
			}
			return strconv.Itoa(count)
		}
	}
	out := runCommand("bash", "-c", "ps -ax | wc -l")
	if out != "" {
		return strings.TrimSpace(out)
	}
	return "n/a"
}

func getCPUTicks() (user, nice, system, idle, iowait, irq, softirq int64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 8 && fields[0] == "cpu" {
			user, _ = strconv.ParseInt(fields[1], 10, 64)
			nice, _ = strconv.ParseInt(fields[2], 10, 64)
			system, _ = strconv.ParseInt(fields[3], 10, 64)
			idle, _ = strconv.ParseInt(fields[4], 10, 64)
			iowait, _ = strconv.ParseInt(fields[5], 10, 64)
			irq, _ = strconv.ParseInt(fields[6], 10, 64)
			softirq, _ = strconv.ParseInt(fields[7], 10, 64)
			return
		}
	}
	return 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid format")
}

func getCPUUsage() string {
	if runtime.GOOS == "linux" {
		u1, n1, s1, id1, io1, ir1, so1, err1 := getCPUTicks()
		if err1 != nil {
			return "n/a"
		}
		time.Sleep(50 * time.Millisecond)
		u2, n2, s2, id2, io2, ir2, so2, err2 := getCPUTicks()
		if err2 != nil {
			return "n/a"
		}

		idle1 := id1 + io1
		idle2 := id2 + io2

		nonIdle1 := u1 + n1 + s1 + ir1 + so1
		nonIdle2 := u2 + n2 + s2 + ir2 + so2

		total1 := idle1 + nonIdle1
		total2 := idle2 + nonIdle2

		totalDiff := total2 - total1
		idleDiff := idle2 - idle1

		if totalDiff > 0 {
			pct := (totalDiff - idleDiff) * 100 / totalDiff
			return fmt.Sprintf("%d%%", pct)
		}
	} else if runtime.GOOS == "darwin" {
		out := runCommand("bash", "-c", "ps -A -o %cpu | awk '{s+=$1} END {print s}'")
		if out != "" {
			if val, err := strconv.ParseFloat(out, 64); err == nil {
				cores := runtime.NumCPU()
				if cores > 0 {
					pct := int(val / float64(cores))
					if pct > 100 {
						pct = 100
					}
					return fmt.Sprintf("%d%%", pct)
				}
			}
		}
	}
	return "n/a"
}

func getCPUTemp() string {
	if runtime.GOOS == "linux" {
		for _, zone := range []string{"thermal_zone0", "thermal_zone1", "thermal_zone2"} {
			data, err := os.ReadFile("/sys/class/thermal/" + zone + "/temp")
			if err == nil {
				tempStr := strings.TrimSpace(string(data))
				if tempVal, err := strconv.ParseFloat(tempStr, 64); err == nil {
					return fmt.Sprintf("%.1f°C", tempVal/1000.0)
				}
			}
		}
		for i := 0; i < 5; i++ {
			for j := 1; j <= 3; j++ {
				path := fmt.Sprintf("/sys/class/hwmon/hwmon%d/temp%d_input", i, j)
				data, err := os.ReadFile(path)
				if err == nil {
					tempStr := strings.TrimSpace(string(data))
					if tempVal, err := strconv.ParseFloat(tempStr, 64); err == nil {
						return fmt.Sprintf("%.1f°C", tempVal/1000.0)
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		out := runCommand("sysctl", "-n", "machdep.xcpm.cpu_thermal_level")
		if out != "" {
			return out + " (level)"
		}
	}
	return "n/a"
}
