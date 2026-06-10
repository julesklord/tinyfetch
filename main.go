package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

var (
    noAscii = flag.Bool("no-ascii", false, "Disable ASCII art")
    theme   = flag.String("theme", "default", "Color theme: default, mono, solar")
)

type Colors struct{
    Bold, Dim, Reset, Cyan, Green, Yellow string
}

func colorsFor(name string) Colors {
    switch name {
    case "mono":
        return Colors{"\x1b[1m","\x1b[2m","\x1b[0m","\x1b[37m","\x1b[37m","\x1b[37m"}
    case "solar":
        return Colors{"\x1b[1m","\x1b[2m","\x1b[0m","\x1b[36m","\x1b[32m","\x1b[33m"}
    default:
        return Colors{"\x1b[1m","\x1b[2m","\x1b[0m","\x1b[36m","\x1b[32m","\x1b[33m"}
    }
}

func runCmd(name string, args ...string) string {
    cmd := exec.Command(name, args...)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return ""
    }
    return strings.TrimSpace(out.String())
}

func readFirstLine(path string) string {
    f, err := os.Open(path)
    if err != nil { return "" }
    defer f.Close()
    r := bufio.NewReader(f)
    s, _ := r.ReadString('\n')
    return strings.TrimSpace(s)
}

func hostname() string {
    if h := runCmd("hostname"); h != "" { return h }
    if h := os.Getenv("HOSTNAME"); h != "" { return h }
    return "unknown"
}

func osPretty() string {
    if out := runCmd("lsb_release", "-ds"); out != "" { return strings.Trim(out, "\"") }
    s := readFirstLine("/etc/os-release")
    for _, line := range strings.Split(s, "\n") {
        if strings.HasPrefix(line, "PRETTY_NAME=") {
            return strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
        }
    }
    return runCmd("uname", "-s")
}

func kernel() string { if k := runCmd("uname", "-r"); k != "" { return k }; return "unknown" }

func uptimePretty() string {
    if data := readFirstLine("/proc/uptime"); data != "" {
        parts := strings.Fields(data)
        if len(parts) > 0 {
            secf, _ := strconv.ParseFloat(parts[0], 64)
            d := time.Duration(secf) * time.Second
            h := int(d.Hours())
            m := int(d.Minutes()) % 60
            return fmt.Sprintf("%dh %dm", h, m)
        }
    }
    if out := runCmd("uptime", "-p"); out != "" { return out }
    return "n/a"
}

func shellName() string { s := os.Getenv("SHELL"); if s=="" { return "-" }; return filepath.Base(s) }

func cpuName() string {
    if s := readFirstLine("/proc/cpuinfo"); s != "" {
        // search for model name
        f, _ := os.Open("/proc/cpuinfo")
        defer f.Close()
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            line := scanner.Text()
            if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Processor") {
                parts := strings.SplitN(line, ":", 2)
                if len(parts) == 2 { return strings.TrimSpace(parts[1]) }
            }
        }
    }
    if out := runCmd("uname", "-p"); out != "" { return out }
    return "unknown"
}

func memoryInfo() (string,string) {
    // returns usage percent and human total (MB/GB)
    f, err := os.Open("/proc/meminfo")
    if err != nil { return "n/a", "n/a" }
    defer f.Close()
    var total, avail int64
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "MemTotal:") {
            fmt.Sscanf(line, "MemTotal: %d kB", &total)
        }
        if strings.HasPrefix(line, "MemAvailable:") {
            fmt.Sscanf(line, "MemAvailable: %d kB", &avail)
        }
    }
    if total == 0 { return "n/a", "n/a" }
    used := total - avail
    pct := int(float64(used)/float64(total)*100 + 0.5)
    human := fmt.Sprintf("%dMB", total/1024)
    if total/1024/1024 > 0 { human = fmt.Sprintf("%dGB", total/1024/1024) }
    return fmt.Sprintf("%d%%", pct), human
}

func diskInfo() string {
    out := runCmd("df", "-h", "--output=source,size,used,avail,pcent,target", "/")
    if out=="" { return "n/a" }
    lines := strings.Split(out, "\n")
    if len(lines) >= 2 {
        parts := strings.Fields(lines[1])
        if len(parts) >= 6 {
            return fmt.Sprintf("%s %s %s (%s avail)", parts[0], parts[1], parts[4], parts[3])
        }
    }
    return "n/a"
}

func packageCount() string {
    if runCmd("which", "dpkg")=="" {
        if runCmd("which", "rpm")=="" {
            return "-"
        }
    }
    if runCmd("which", "dpkg")!="" {
        out := runCmd("dpkg", "-l")
        if out=="" { return "-" }
        // rough count
        return strconv.Itoa(len(strings.Split(out, "\n"))-5)
    }
    if runCmd("which", "rpm")!="" {
        out := runCmd("rpm","-qa")
        if out=="" { return "-" }
        return strconv.Itoa(len(strings.Split(out, "\n")))
    }
    return "-"
}

func resolution() string {
    if runCmd("which","xdpyinfo")!="" {
        out := runCmd("xdpyinfo")
        for _, line := range strings.Split(out, "\n") {
            if strings.Contains(line, "dimensions:") {
                return strings.TrimSpace(strings.SplitN(line, ":",2)[1])
            }
        }
    }
    return "-"
}

func gitBranch() string {
    if runCmd("which","git")=="" { return "-" }
    if runCmd("git","rev-parse","--is-inside-work-tree") == "true" {
        if b := runCmd("git","rev-parse","--abbrev-ref","HEAD"); b != "" { return b }
    }
    return "-"
}

func gpuInfo() string {
    if runCmd("which","lspci")=="" { return "-" }
    out := runCmd("lspci")
    var res []string
    for _, line := range strings.Split(out, "\n") {
        lower := strings.ToLower(line)
        if strings.Contains(lower, "vga") || strings.Contains(lower, "3d") || strings.Contains(lower, "display") {
            // take description after the device id
            parts := strings.SplitN(line, ":", 2)
            if len(parts) == 2 { res = append(res, strings.TrimSpace(parts[1])) }
        }
    }
    if len(res)>0 { return strings.Join(res, "; ") }
    return "-"
}

func batteryInfo() string {
    // check /sys/class/power_supply for BAT*
    base := "/sys/class/power_supply"
    entries, err := os.ReadDir(base)
    if err != nil { return "-" }
    for _, e := range entries {
        if strings.HasPrefix(e.Name(), "BAT") {
            cap := readFirstLine(filepath.Join(base, e.Name(), "capacity"))
            stat := readFirstLine(filepath.Join(base, e.Name(), "status"))
            if cap == "" { cap = "?" }
            if stat == "" { stat = "unknown" }
            return fmt.Sprintf("%s%% (%s)", cap, stat)
        }
    }
    return "-"
}

func desktopEnv() string {
    if v := os.Getenv("XDG_CURRENT_DESKTOP"); v != "" { return v }
    if v := os.Getenv("DESKTOP_SESSION"); v != "" { return v }
    if v := os.Getenv("WAYLAND_DISPLAY"); v != "" { return "Wayland" }
    return "-"
}

func printKV(c Colors, k, v string) {
    fmt.Printf("%s%-14s %s%s\n", c.Dim, k+":", c.Bold, v)
}

func main() {
    flag.Parse()
    c := colorsFor(*theme)

    // gather
    host := hostname()
    osname := osPretty()
    kern := kernel()
    up := uptimePretty()
    shell := shellName()
    cpu := cpuName()
    memPct, memHuman := memoryInfo()
    disk := diskInfo()
    pkgs := packageCount()
    res := resolution()
    term := os.Getenv("TERM")
    branch := gitBranch()
    gpu := gpuInfo()
    bat := batteryInfo()
    de := desktopEnv()

    right := [][]string{
        {"Host", host},
        {"OS", osname},
        {"Kernel", kern},
        {"Uptime", up},
        {"Shell", shell},
        {"CPU", cpu},
        {"Memory", fmt.Sprintf("%s (%s)", memPct, memHuman)},
        {"Disk", disk},
        {"Packages", pkgs},
        {"Resolution", res},
        {"Terminal", term},
        {"Git branch", branch},
        {"GPU", gpu},
        {"Battery", bat},
        {"DE", de},
    }

    if *noAscii {
        for _, kv := range right { printKV(c, kv[0], kv[1]) }
        return
    }

    ascii := []string{
        `  __  __ _ `,
        ` |  \/  (_)`,
        ` | |\/| | |`,
        ` | |  | | |`,
        ` |_|  |_|_|`,
    }

    // ensure ascii lines >= right lines
    if len(ascii) < len(right) {
        for i := len(ascii); i < len(right); i++ { ascii = append(ascii, "") }
    }

    for i := 0; i < len(ascii); i++ {
        left := ascii[i]
        if i < len(right) {
            kv := right[i]
            fmt.Printf("%s%s %s %-12s %s%s\n", c.Cyan, left, c.Reset, kv[0]+":", c.Bold, kv[1])
        } else {
            fmt.Printf("%s%s%s\n", c.Cyan, left, c.Reset)
        }
    }
}
