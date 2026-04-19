package version_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/version"
)

func TestGet_ReturnsDefaults(t *testing.T) {
	info := version.Get()
	if info.Version == "" {
		t.Error("expected non-empty Version")
	}
	if info.Commit == "" {
		t.Error("expected non-empty Commit")
	}
	if info.BuildDate == "" {
		t.Error("expected non-empty BuildDate")
	}
}

func TestString_ContainsVersion(t *testing.T) {
	version.Version = "1.2.3"
	info := version.Get()
	s := info.String()
	if !strings.Contains(s, "1.2.3") {
		t.Errorf("String() = %q, want it to contain version", s)
	}
	if !strings.Contains(s, "portwatch") {
		t.Errorf("String() = %q, want it to contain 'portwatch'", s)
	}
}

func TestString_ContainsCommitAndDate(t *testing.T) {
	version.Commit = "abc1234"
	version.BuildDate = "2024-01-15"
	info := version.Get()
	s := info.String()
	if !strings.Contains(s, "abc1234") {
		t.Errorf("String() = %q, want commit hash", s)
	}
	if !strings.Contains(s, "2024-01-15") {
		t.Errorf("String() = %q, want build date", s)
	}
}

func TestShort_ReturnsVersionOnly(t *testing.T) {
	version.Version = "2.0.0"
	info := version.Get()
	if info.Short() != "2.0.0" {
		t.Errorf("Short() = %q, want %q", info.Short(), "2.0.0")
	}
}
