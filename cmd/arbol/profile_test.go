package main

import (
	"testing"
)

func BenchmarkGatherInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gatherInfo("../../plugins")
	}
}
