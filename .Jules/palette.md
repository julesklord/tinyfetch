## 2024-06-30 - Invalid CLI Flag Feedback
**Learning:** The CLI tool was continuing execution and outputting mangled data when provided with unknown flags, causing user confusion. Providing immediate, clear error feedback with usage instructions improves developer experience.
**Action:** Always validate all input flags and exit early with helpful error messages and usage instructions for unknown arguments.

## 2024-07-01 - Changed default bar style
**Learning:** Braille characters in terminal can be hard to read for some users or terminal emulators without proper font support.
**Action:** Change default bar style to `BarStyleBlock`.
