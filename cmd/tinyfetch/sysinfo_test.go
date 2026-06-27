package main

import (
	"testing"
)

func TestRunCommand_Success(t *testing.T) {
	out := runCommand("echo", "hello")
	if out != "hello" {
		t.Errorf("Expected 'hello', got %q", out)
	}
}

func TestRunCommand_ErrorNotFound(t *testing.T) {
	out := runCommand("this-command-does-not-exist-12345")
	if out != "" {
		t.Errorf("Expected empty string, got %q", out)
	}
}

func TestRunCommand_ErrorExitCode(t *testing.T) {
	// A command that exists but fails with non-zero exit code
	out := runCommand("sh", "-c", "exit 1")
	if out != "" {
		t.Errorf("Expected empty string on error, got %q", out)
	}
}
