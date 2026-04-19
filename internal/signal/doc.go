// Package signal wraps OS signal handling for portwatch.
//
// It provides a Handler that translates SIGINT / SIGTERM (or any
// caller-supplied signal) into context cancellation, enabling clean
// shutdown of the scan loop and background goroutines.
//
// Basic usage:
//
//	h := signal.New(os.Stderr)
//	ctx, stop := h.Notify(context.Background())
//	defer stop()
//	signal.Wait(ctx)
package signal
