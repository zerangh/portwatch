package lock

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func tmpState(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "state.json")
}

func TestNew_LockPathDerived(t *testing.T) {
	l := New("/var/lib/portwatch/state.json")
	want := "/var/lib/portwatch/.state.json.lock"
	if l.Path() != want {
		t.Fatalf("got %q, want %q", l.Path(), want)
	}
}

func TestAcquireAndRelease(t *testing.T) {
	l := New(tmpState(t))
	if err := l.Acquire(); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	defer func() {
		if err := l.Release(); err != nil {
			t.Errorf("Release: %v", err)
		}
	}()

	data, err := os.ReadFile(l.Path())
	if err != nil {
		t.Fatalf("lock file missing: %v", err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil || pid != os.Getpid() {
		t.Fatalf("unexpected lock content %q", data)
	}
}

func TestAcquire_StaleLockOverwritten(t *testing.T) {
	l := New(tmpState(t))
	// Write a lock file with a PID that cannot exist (0).
	if err := os.WriteFile(l.Path(), []byte("0"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := l.Acquire(); err != nil {
		t.Fatalf("expected stale lock to be overwritten, got: %v", err)
	}
	_ = l.Release()
}

func TestAcquire_LiveLockReturnsError(t *testing.T) {
	l := New(tmpState(t))
	// Write a lock file with our own PID — simulates another process.
	pid := strconv.Itoa(os.Getpid())
	if err := os.WriteFile(l.Path(), []byte(pid), 0o600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(l.Path())

	err := l.Acquire()
	if err == nil {
		_ = l.Release()
		t.Fatal("expected ErrLocked, got nil")
	}
}

func TestRelease_MissingFile(t *testing.T) {
	l := New(tmpState(t))
	// Releasing a lock that was never acquired should not panic.
	if err := l.Release(); err == nil {
		t.Fatal("expected error releasing non-existent lock")
	}
}
