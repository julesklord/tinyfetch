package main

import (
	"testing"
)

func BenchmarkGetProcesses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getProcesses()
	}
}
