// Package portwatch provides the high-level scan pipeline and runner that
// ties together scanning, diffing, alerting, and state persistence.
//
// # Pipeline
//
// A Pipeline performs a single watch cycle:
//
//  1. Scan the configured port range via [scanner.Scanner].
//  2. Compare results against the previously saved state using [state.Diff].
//  3. Emit alerts for any changes via [alert.Alerter].
//  4. Persist the new snapshot to disk with [state.Save].
//
// # Runner
//
// A Runner wraps a Pipeline and drives it on a fixed [time.Duration] interval.
// It runs the pipeline immediately on start, then on every tick until the
// context is cancelled or an optional MaxRuns limit is reached.
//
// Example:
//
//	cfg, _ := config.Load("portwatch.toml")
//	p, _ := portwatch.NewPipeline(cfg, "/var/lib/portwatch/state.json", os.Stdout)
//	r, _ := portwatch.NewRunner(p, portwatch.RunnerConfig{Interval: 60 * time.Second}, os.Stderr)
//	r.Start(ctx)
package portwatch
