//go:build !linux

package main

import (
	"strings"
)

func getProcesses() string {
	out := runCommand("bash", "-c", "ps -ax | wc -l")
	if out != "" {
		return strings.TrimSpace(out)
	}
	return "n/a"
}
