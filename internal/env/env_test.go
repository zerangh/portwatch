package env_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/env"
)

func TestLoad_EmptyEnv(t *testing.T) {
	clearEnv(t)
	v := env.Load()
	if v.IsSet() {
		t.Fatal("expected no values set")
	}
}

func TestLoad_PortRange(t *testing.T) {
	clearEnv(t)
	t.Setenv(env.EnvPortRange, "1-1024")
	v := env.Load()
	if v.PortRange != "1-1024" {
		t.Fatalf("got %q, want %q", v.PortRange, "1-1024")
	}
}

func TestLoad_Concurrency(t *testing.T) {
	clearEnv(t)
	t.Setenv(env.EnvConcurrency, "8")
	v := env.Load()
	if v.Concurrency != 8 {
		t.Fatalf("got %d, want 8", v.Concurrency)
	}
}

func TestLoad_InvalidConcurrency(t *testing.T) {
	clearEnv(t)
	t.Setenv(env.EnvConcurrency, "bad")
	v := env.Load()
	if v.Concurrency != 0 {
		t.Fatalf("expected 0 for invalid concurrency, got %d", v.Concurrency)
	}
}

func TestLoad_Interval(t *testing.T) {
	clearEnv(t)
	t.Setenv(env.EnvInterval, "30s")
	v := env.Load()
	if v.Interval != 30*time.Second {
		t.Fatalf("got %v, want 30s", v.Interval)
	}
}

func TestLoad_Webhook(t *testing.T) {
	clearEnv(t)
	t.Setenv(env.EnvAlertWebhook, "https://example.com/hook")
	v := env.Load()
	if v.Webhook != "https://example.com/hook" {
		t.Fatalf("got %q", v.Webhook)
	}
	if !v.IsSet() {
		t.Fatal("expected IsSet true")
	}
}

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		env.EnvPortRange, env.EnvConcurrency,
		env.EnvInterval, env.EnvStatePath, env.EnvAlertWebhook,
	} {
		t.Setenv(k, "")
	}
}
