package main

import (
	"testing"
)

func TestRuneWidth(t *testing.T) {
	tests := []struct {
		r    rune
		want int
	}{
		{'A', 1},
		{'\u200d', 0}, // Zero-width joiner
		{'界', 2},      // CJK character
		{'🌸', 2},      // Emoji
	}

	for _, tt := range tests {
		got := runeWidth(tt.r)
		if got != tt.want {
			t.Errorf("runeWidth(%q) = %d; want %d", tt.r, got, tt.want)
		}
	}
}

func TestVisualLength(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"ASCII plain", "hello", 5},
		{"ANSI color strip", "\033[01;34mHost:\033[0m   myhost", 14},
		{"CJK characters", "hello世界", 9},               // 5 (hello) + 4 (世界)
		{"Nerd Fonts / Emojis", "Host:  Spotify", 15}, // 6 (Host: ) + 1 () + 8 ( Spotify)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := visualLength(tt.s)
			if got != tt.want {
				t.Errorf("visualLength(%q) = %d; want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestTruncateANSI(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		limit int
		want  string
	}{
		{
			name:  "no truncation needed",
			s:     "hello",
			limit: 10,
			want:  "hello",
		},
		{
			name:  "basic truncation",
			s:     "hello world",
			limit: 8,
			want:  "hello w…\033[0m", // limit is 8, targetLen is 7, "hello w" (7) + "…" (1) + reset
		},
		{
			name:  "truncation with ANSI colors",
			s:     "\033[01;34mhello world\033[0m",
			limit: 8,
			want:  "\033[01;34mhello w\033[0m…\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateANSI(tt.s, tt.limit)
			if got != tt.want {
				t.Errorf("truncateANSI(%q, %d) = %q; want %q", tt.s, tt.limit, got, tt.want)
			}
		})
	}
}
