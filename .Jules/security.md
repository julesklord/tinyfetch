## 2026-06-30 - Improve ANSI escape sequence stripping
**Learning:** The previous ANSI sequence stripper only expected 'm' as the terminator for escapes, leading to vulnerabilities or improper rendering when given sequences ending in other characters. A proper ANSI parser should expect any character in the range 0x40-0x7E for CSI sequences.
**Action:** Use explicit CSI tracking state to correctly handle all ANSI escape sequence terminators in both 'stripANSI' and 'truncateANSI'.
