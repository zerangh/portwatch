package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestMainBinaryBuilds verifies the binary compiles without errors.
func TestMainBinaryBuilds(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping build test on windows")
	}

	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "portwatch")

	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if _, err := os.Stat(binary); os.IsNotExist(err) {
		t.Fatal("binary not produced after build")
	}
}

// TestMainHelp runs the binary with -help and expects a clean exit.
func TestMainHelp(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}

	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "portwatch")

	build := exec.Command("go", "build", "-o", binary, ".")
	if err := build.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	cmd := exec.Command(binary, "-help")
	// -help exits with code 2 for flag package, that's expected
	out, _ := cmd.CombinedOutput()
	if len(out) == 0 {
		t.Error("expected usage output from -help flag")
	}
}
