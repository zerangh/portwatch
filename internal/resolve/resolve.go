// Package resolve provides hostname-to-IP resolution utilities for portwatch.
package resolve

import (
	"fmt"
	"net"
	"time"
)

// Result holds the resolved addresses for a hostname.
type Result struct {
	Host      string
	Addresses []string
	ResolvedAt time.Time
}

// Resolver resolves hostnames to IP addresses.
type Resolver struct {
	timeout time.Duration
}

// New returns a Resolver with the given timeout.
func New(timeout time.Duration) *Resolver {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Resolver{timeout: timeout}
}

// Resolve looks up the IP addresses for the given host.
func (r *Resolver) Resolve(host string) (*Result, error) {
	if host == "" {
		return nil, fmt.Errorf("resolve: host must not be empty")
	}

	// If already an IP, return as-is.
	if ip := net.ParseIP(host); ip != nil {
		return &Result{
			Host:       host,
			Addresses:  []string{ip.String()},
			ResolvedAt: time.Now(),
		}, nil
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, fmt.Errorf("resolve: lookup %q: %w", host, err)
	}

	return &Result{
		Host:       host,
		Addresses:  addrs,
		ResolvedAt: time.Now(),
	}, nil
}

// Primary returns the first resolved address, or an error if none exist.
func (res *Result) Primary() (string, error) {
	if len(res.Addresses) == 0 {
		return "", fmt.Errorf("resolve: no addresses for %q", res.Host)
	}
	return res.Addresses[0], nil
}
