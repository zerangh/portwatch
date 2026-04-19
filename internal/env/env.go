// Package env provides helpers for reading portwatch configuration
// from environment variables, complementing file-based config.
package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EnvPortRange    = "PORTWATCH_PORT_RANGE"
	EnvConcurrency  = "PORTWATCH_CONCURRENCY"
	EnvInterval     = "PORTWATCH_INTERVAL"
	EnvStatePath    = "PORTWATCH_STATE_PATH"
	EnvAlertWebhook = "PORTWATCH_ALERT_WEBHOOK"
)

// Values holds environment-sourced configuration overrides.
type Values struct {
	PortRange   string
	Concurrency int
	Interval    time.Duration
	StatePath   string
	Webhook     string
}

// Load reads known environment variables and returns a Values struct.
// Fields retain zero values when the corresponding variable is unset or invalid.
func Load() Values {
	v := Values{
		PortRange: os.Getenv(EnvPortRange),
		StatePath: os.Getenv(EnvStatePath),
		Webhook:   os.Getenv(EnvAlertWebhook),
	}

	if s := os.Getenv(EnvConcurrency); s != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && n > 0 {
			v.Concurrency = n
		}
	}

	if s := os.Getenv(EnvInterval); s != "" {
		if d, err := time.ParseDuration(strings.TrimSpace(s)); err == nil && d > 0 {
			v.Interval = d
		}
	}

	return v
}

// IsSet reports whether any environment override was detected.
func (v Values) IsSet() bool {
	return v.PortRange != "" || v.Concurrency != 0 || v.Interval != 0 ||
		v.StatePath != "" || v.Webhook != ""
}
