package portwatch

import (
	"bytes"
	"testing"
	"time"
)

func TestScanTimingMiddleware_RecordsDuration(t *testing.T) {
	var buf bytes.Buffer
	h := ScanTimingMiddleware(&buf)

	now := time.Now()
	h(EventScanBegin, now)
	h(EventScanEnd, now.Add(120*time.Millisecond))

	if !bytes.Contains(buf.Bytes(), []byte("[timing]")) {
		t.Errorf("expected timing line, got: %s", buf.String())
	}
}

func TestScanTimingMiddleware_NoBegin_Silent(t *testing.T) {
	var buf bytes.Buffer
	h := ScanTimingMiddleware(&buf)
	h(EventScanEnd, time.Now())
	if buf.Len() != 0 {
		t.Errorf("expected no output without prior scan_begin, got: %s", buf.String())
	}
}

func TestScanTimingMiddleware_NilWriterUsesStdout(t *testing.T) {
	h := ScanTimingMiddleware(nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestUptimeMiddleware_PrintsUptime(t *testing.T) {
	var buf bytes.Buffer
	h := UptimeMiddleware(&buf)

	now := time.Now()
	h(EventReady, now)
	h(EventStopping, now.Add(5*time.Second))

	if !bytes.Contains(buf.Bytes(), []byte("[uptime]")) {
		t.Errorf("expected uptime line, got: %s", buf.String())
	}
}

func TestUptimeMiddleware_NoReady_Silent(t *testing.T) {
	var buf bytes.Buffer
	h := UptimeMiddleware(&buf)
	h(EventStopping, time.Now())
	if buf.Len() != 0 {
		t.Errorf("expected no output without prior ready event, got: %s", buf.String())
	}
}

func TestUptimeMiddleware_NilWriterUsesStdout(t *testing.T) {
	h := UptimeMiddleware(nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestOtherEvents_Ignored(t *testing.T) {
	var buf bytes.Buffer
	th := ScanTimingMiddleware(&buf)
	uh := UptimeMiddleware(&buf)
	now := time.Now()
	for _, ev := range []LifecycleEvent{EventStarting, EventStopped} {
		th(ev, now)
		uh(ev, now)
	}
	if buf.Len() != 0 {
		t.Errorf("unexpected output for unrelated events: %s", buf.String())
	}
}
