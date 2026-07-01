package main

import (
	"strings"
	"unicode/utf8"
)

type BarStyle int

const (
	BarStyleBlock BarStyle = iota
	BarStyleBraille
	BarStyleGradient
	BarStyleDot
)

var currentBarStyle = BarStyleBlock

func SetBarStyle(style BarStyle) {
	currentBarStyle = style
}

func runeWidth(r rune) int {
	// Zero-width space, joiners, control chars, variation selectors
	if r == '\u200d' || r == '\u200c' || (r >= '\ufe00' && r <= '\ufe0f') {
		return 0
	}
	// Combining diacritical marks
	if r >= 0x0300 && r <= 0x036F {
		return 0
	}
	// Wide ranges (2 columns)
	// Emojis / Pictographs in SMP (Plane 1): U+1F000 to U+ to U+1FAFF
	if r >= 0x1F000 && r <= 0x1FAFF {
		return 2
	}
	// Miscellaneous Symbols and Pictographs, Emoticons, Ornamental Dingbats, etc. in BMP
	if r >= 0x2600 && r <= 0x27BF {
		return 2
	}
	// Braille patterns
	if r >= 0x2800 && r <= 0x28FF {
		return 1
	}
	// CJK ranges
	if (r >= 0x2E80 && r <= 0x2FDF) || // CJK Radicals
		(r >= 0x3000 && r <= 0x9FFF) || // Hiragana, Katakana, CJK Unified Ideographs
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility
		(r >= 0xFF01 && r <= 0xFF60) || // Fullwidth Forms
		(r >= 0xFFE0 && r <= 0xFFE6) {
		return 2
	}
	// Default
	return 1
}

func visualLength(s string) int {
	raw := stripANSI(s)
	length := 0
	for _, r := range raw {
		length += runeWidth(r)
	}
	return length
}

func truncateANSI(s string, limit int) string {
	if visualLength(s) <= limit {
		return s
	}

	var builder strings.Builder
	visualLen := 0
	inEscape := false
	isCSI := false
	restoreCode := "\033[0m"
	targetLen := limit - 1
	if targetLen < 0 {
		targetLen = 0
	}

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			isCSI = false
			builder.WriteByte(s[i])
			continue
		}
		if inEscape {
			builder.WriteByte(s[i])
			if !isCSI {
				if s[i] == '[' {
					isCSI = true
				} else {
					inEscape = false
				}
				continue
			}
			if s[i] >= 0x40 && s[i] <= 0x7E {
				inEscape = false
				isCSI = false
			}
			continue
		}

		if visualLen < targetLen {
			r, size := utf8.DecodeRuneInString(s[i:])
			w := runeWidth(r)
			if visualLen+w <= targetLen {
				builder.WriteRune(r)
				visualLen += w
			} else {
				visualLen = targetLen
			}
			i += size - 1
		}
	}
	builder.WriteString("…")
	builder.WriteString(restoreCode)
	return builder.String()
}

func stripANSI(s string) string {
	var builder strings.Builder
	inEscape := false
	isCSI := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			isCSI = false
			continue
		}
		if inEscape {
			if !isCSI {
				if s[i] == '[' {
					isCSI = true
				} else {
					inEscape = false
				}
				continue
			}
			if s[i] >= 0x40 && s[i] <= 0x7E {
				inEscape = false
				isCSI = false
			}
			continue
		}
		builder.WriteByte(s[i])
	}
	return builder.String()
}

func getBar(pct int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}

	color := "\033[01;32m" // Green
	if pct > 80 {
		color = "\033[01;31m" // Red
	} else if pct > 50 {
		color = "\033[01;33m" // Yellow
	}
	restore := "\033[0m"
	gray := "\033[00;37m"

	switch currentBarStyle {
	case BarStyleBraille:
		return getBrailleBar(pct, color, gray, restore)
	case BarStyleGradient:
		return getGradientBar(pct, restore, gray)
	case BarStyleDot:
		return getDotBar(pct, color, gray, restore)
	default:
		return getBlockBar(pct, color, gray, restore)
	}
}

func getBlockBar(pct int, color, gray, restore string) string {
	filled := pct / 10
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled

	var sb strings.Builder
	sb.WriteString(color)
	for i := 0; i < filled; i++ {
		sb.WriteString("█")
	}
	sb.WriteString(restore + gray)
	for i := 0; i < empty; i++ {
		sb.WriteString("░")
	}
	sb.WriteString(restore)
	return sb.String()
}

func getBrailleBar(pct int, color, gray, restore string) string {
	// Braille patterns for 8x resolution (each char = 8 segments)
	// We use 10 chars = 80 segments for 0-100%
	braillePatterns := []string{
		"⠀", // 0/8
		"⠁", // 1/8
		"⠃", // 2/8
		"⠇", // 3/8
		"⠏", // 4/8
		"⠟", // 5/8
		"⠿", // 6/8 (close enough)
		"⠿", // 7/8
		"⠿", // 8/8
	}

	totalSegments := 80 // 10 chars * 8 segments
	filledSegments := pct * totalSegments / 100
	fullChars := filledSegments / 8
	partialSegments := filledSegments % 8

	var sb strings.Builder
	sb.WriteString(color)
	for i := 0; i < fullChars; i++ {
		sb.WriteString("⠿")
	}
	if fullChars < 10 {
		sb.WriteString(braillePatterns[partialSegments])
		sb.WriteString(restore + gray)
		for i := fullChars + 1; i < 10; i++ {
			sb.WriteString("⠀")
		}
	} else {
		sb.WriteString(restore + gray)
	}
	sb.WriteString(restore)
	return sb.String()
}

func getGradientBar(pct int, restore, gray string) string {
	// 4-color gradient: green -> yellow -> orange -> red
	gradientChars := []string{"░", "▒", "▓", "█"}
	gradientColors := []string{
		"\033[38;2;0;255;0m",     // Green
		"\033[38;2;170;255;0m",   // Yellow-green
		"\033[38;2;255;170;0m",   // Orange
		"\033[38;2;255;0;0m",     // Red
	}

	filled := pct / 10
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled

	var sb strings.Builder
	for i := 0; i < filled; i++ {
		colorIdx := i * 4 / 10
		if colorIdx > 3 {
			colorIdx = 3
		}
		sb.WriteString(gradientColors[colorIdx])
		sb.WriteString(gradientChars[3])
	}
	sb.WriteString(restore + gray)
	for i := 0; i < empty; i++ {
		sb.WriteString(gradientChars[0])
	}
	sb.WriteString(restore)
	return sb.String()
}

func getDotBar(pct int, color, gray, restore string) string {
	filled := pct / 10
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled

	var sb strings.Builder
	sb.WriteString(color)
	for i := 0; i < filled; i++ {
		sb.WriteString("●")
	}
	sb.WriteString(restore + gray)
	for i := 0; i < empty; i++ {
		sb.WriteString("○")
	}
	sb.WriteString(restore)
	return sb.String()
}
