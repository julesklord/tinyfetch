//go:build linux

package main

import (
	"strconv"
	"strings"
	"syscall"
)

func getProcesses() string {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err == nil {
		return strconv.FormatUint(uint64(info.Procs), 10)
	}
	out := runCommand("bash", "-c", "ps -ax | wc -l")
	if out != "" {
		return strings.TrimSpace(out)
	}
	return "n/a"
}
