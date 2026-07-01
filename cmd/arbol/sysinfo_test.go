package main

import (
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	out := runCommand("echo", "hello")
	if out != "hello" {
		t.Errorf("Expected 'hello', got '%s'", out)
	}

	outFail := runCommand("this-command-does-not-exist-12345")
	if outFail != "" {
		t.Errorf("Expected empty string on failure, got '%s'", outFail)
	}
}

func TestRunCommand_ExitCode(t *testing.T) {
	// A command that exists but fails with non-zero exit code
	out := runCommand("sh", "-c", "exit 1")
	if out != "" {
		t.Errorf("Expected empty string for non-zero exit code, got %q", out)
	}
}

func TestRunCommandWithTimeout(t *testing.T) {
	out := runCommandWithTimeout(1*time.Second, "echo", "hello")
	if out != "hello" {
		t.Errorf("Expected 'hello', got '%s'", out)
	}

	outFail := runCommandWithTimeout(1*time.Second, "this-command-does-not-exist-12345")
	if outFail != "" {
		t.Errorf("Expected empty string on failure, got '%s'", outFail)
	}

	// Test timeout scenario
	outTimeout := runCommandWithTimeout(10*time.Millisecond, "sleep", "1")
	if outTimeout != "" {
		t.Errorf("Expected empty string on timeout, got '%s'", outTimeout)
	}
}

func TestSysInfoFunctionsSanity(t *testing.T) {
	// These functions execute OS-specific commands.
	// The exact output varies by system, but they should not panic
	// and should generally return non-empty strings (even if it's "n/a" or "Unknown").

	funcs := []struct {
		name string
		fn   func() string
	}{
		{"getOSName", getOSName},
		{"getDistroID", getDistroID},
		{"getUptime", getUptime},
		{"getCPU", getCPU},
		{"getMemory", getMemory},
		{"getDisk", getDisk},
	}

	for _, tc := range funcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.fn()
			if res == "" {
				t.Errorf("%s() returned empty string", tc.name)
			}
		})
	}
}
