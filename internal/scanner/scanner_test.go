package scanner

import (
	"net"
	"testing"
	"time"
)

func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_FindsOpenPort(t *testing.T) {
	port, stop := startTestServer(t)
	defer stop()

	s := New("127.0.0.1", 500*time.Millisecond)
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0].Number != port {
		t.Errorf("expected open port %d, got %v", port, ports)
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := New("127.0.0.1", 100*time.Millisecond)
	_, err := s.Scan(500, 100)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestPortString(t *testing.T) {
	p := Port{Protocol: "tcp", Number: 8080}
	if p.String() != "tcp/8080" {
		t.Errorf("unexpected string: %s", p.String())
	}
}
