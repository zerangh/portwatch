// Package watchdog provides a liveness monitor for port scan cycles.
//
// A Watchdog is created with a maximum acceptable age for scan completions.
// After each successful scan the caller should call Ping to reset the timer.
// Check can be called periodically to emit a warning when the watchdog
// has not been pinged within the configured interval, indicating that
// the scan loop may have stalled or crashed.
//
// Example:
//
//	wd, _ := watchdog.New(2*time.Minute, os.Stderr)
//	// inside scan loop:
//	wd.Ping()
//	// in a separate goroutine or health endpoint:
//	wd.Check()
package watchdog
