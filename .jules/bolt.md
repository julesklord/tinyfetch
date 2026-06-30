## 2024-06-30 - Parallelizing Mandatory System Sleep Blocks
**Learning:** System information collection routines (like CPU usage sampling) often require mandatory sleep intervals (e.g., a 50ms wait to sample CPU ticks). In a sequential pipeline, these sleeps directly block overall execution.
**Action:** Always look for blocking sleep or high-latency I/O operations in sequential data gathering functions. Wrap them in a goroutine and read from a channel at the end of the pipeline to effectively hide the latency behind parallel execution of other tasks.
