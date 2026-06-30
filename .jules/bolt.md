## 2026-06-28 - Optimize strings.Split parsing loop
**Learning:** Using `strings.Split` in tight loops creates temporary slice allocations which increases GC pressure and slows execution. For large or repetitively processed strings where only sequential access is required, iterative `strings.IndexByte` is significantly faster.
**Action:** Prefer manual string splitting loops using `strings.IndexByte` or `strings.Cut` to avoid unnecessary slice allocations when processing multi-line outputs line-by-line.
