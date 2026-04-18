package config_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourusername/portwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.PortRange != "1-1024" {
		t.Errorf("expected default port range '1-1024', got %q", cfg.PortRange)
	}
	if cfg.Concurrency != 100 {
		t.Errorf("expected default concurrency 100, got %d", cfg.Concurrency)
	}
	if cfg.StateFile == "" {
		t.Error("expected non-empty default state file")
	}
}

func TestLoad_MissingFile_ReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/portwatch.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg.PortRange != "1-1024" {
		t.Errorf("expected default port range, got %q", cfg.PortRange)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	data := map[string]interface{}{
		"port_range":  "80-443",
		"state_file":  "/tmp/state.json",
		"concurrency": 50,
		"alert_email": "ops@example.com",
	}
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := config.Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PortRange != "80-443" {
		t.Errorf("expected port range '80-443', got %q", cfg.PortRange)
	}
	if cfg.AlertEmail != "ops@example.com" {
		t.Errorf("expected alert email, got %q", cfg.AlertEmail)
	}
	if cfg.Concurrency != 50 {
		t.Errorf("expected concurrency 50, got %d", cfg.Concurrency)
	}
}

func TestValidate_InvalidConcurrency(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Concurrency = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for zero concurrency")
	}
}

func TestValidate_EmptyPortRange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.PortRange = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty port range")
	}
}
