// Package resolve provides hostname and IP address resolution for portwatch.
//
// It wraps standard net.LookupHost with a simple Result type that captures
// resolved addresses and the time of resolution. If the input is already a
// valid IP address it is returned directly without a DNS lookup.
//
// Usage:
//
//	r := resolve.New(5 * time.Second)
//	res, err := r.Resolve("example.com")
//	ip, err := res.Primary()
package resolve
