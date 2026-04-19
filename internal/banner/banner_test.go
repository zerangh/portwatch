package banner_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/banner"
)

func TestPrint_ContainsVersion(t *testing.T) {
	var buf bytes.Buffer
	banner.Print(banner.Options{Writer: &buf})
	if !strings.Contains(buf.String(), "0.1.0") {
		t.Errorf("expected version in banner, got:\n%s", buf.String())
	}
}

func TestPrint_ContainsHost(t *testing.T) {
	var buf bytes.Buffer
	banner.Print(banner.Options{
		Host:   "example.com",
		Writer: &buf,
	})
	if !strings.Contains(buf.String(), "example.com") {
		t.Errorf("expected host in banner output")
	}
}

func TestPrint_ContainsPortRange(t *testing.T) {
	var buf bytes.Buffer
	banner.Print(banner.Options{
		PortRange: "1-1024",
		Writer:    &buf,
	})
	if !strings.Contains(buf.String(), "1-1024") {
		t.Errorf("expected port range in banner output")
	}
}

func TestPrint_ContainsInterval(t *testing.T) {
	var buf bytes.Buffer
	banner.Print(banner.Options{
		Interval: 30 * time.Second,
		Writer:   &buf,
	})
	if !strings.Contains(buf.String(), "30s") {
		t.Errorf("expected interval in banner output")
	}
}

func TestPrint_NilWriterUsesStdout(t *testing.T) {
	// Should not panic when Writer is nil.
	banner.Print(banner.Options{
		Host:      "localhost",
		PortRange: "80-443",
		Interval:  time.Minute,
		StatePath: "/tmp/state.json",
	})
}
