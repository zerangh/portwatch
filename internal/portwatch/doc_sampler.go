// Package portwatch provides the core pipeline, runner, and supporting
// primitives used to orchestrate periodic port scans.
//
// # Sampler
//
// Sampler collects scan-duration observations in a fixed-size sliding window
// and exposes rolling statistics useful for performance monitoring:
//
//	s := portwatch.NewSampler(100, os.Stdout)
//
//	start := time.Now()
//	// … run scan …
//	s.Record(time.Since(start))
//
//	fmt.Println("mean:", s.Mean())
//	fmt.Println("p95: ", s.P95())
//
// The window evicts the oldest sample when it reaches capacity, so statistics
// always reflect the most recent observations.
package portwatch
