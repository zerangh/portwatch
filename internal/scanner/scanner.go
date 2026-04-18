package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"
)

// Port represents an open port with its protocol and number.
type Port struct {
	Protocol string
	Number   int
}

func (p Port) String() string {
	return fmt.Sprintf("%s/%d", p.Protocol, p.Number)
}

// Scanner scans for open ports on a host.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner for the given host.
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{Host: host, Timeout: timeout}
}

// Scan checks which ports in the given range are open via TCP.
func (s *Scanner) Scan(startPort, endPort int) ([]Port, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var open []Port

	for port := startPort; port <= endPort; port++ {
		addr := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err == nil {
			conn.Close()
			open = append(open, Port{Protocol: "tcp", Number: port})
		}
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i].Number < open[j].Number
	})

	return open, nil
}
