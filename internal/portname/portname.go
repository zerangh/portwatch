// Package portname maps well-known port numbers to service names.
package portname

import "fmt"

// well-known port-to-service mappings.
var builtins = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Resolver resolves port numbers to human-readable service names.
type Resolver struct {
	extra map[int]string
}

// New returns a Resolver with optional extra mappings merged over the builtins.
func New(extra map[int]string) *Resolver {
	m := make(map[int]string, len(extra))
	for k, v := range extra {
		m[k] = v
	}
	return &Resolver{extra: m}
}

// Lookup returns the service name for port, falling back to "port/<n>".
func (r *Resolver) Lookup(port int) string {
	if r != nil {
		if name, ok := r.extra[port]; ok {
			return name
		}
	}
	if name, ok := builtins[port]; ok {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

// Annotate returns a slice of "port(name)" strings for display.
func (r *Resolver) Annotate(ports []int) []string {
	out := make([]string, len(ports))
	for i, p := range ports {
		out[i] = fmt.Sprintf("%d(%s)", p, r.Lookup(p))
	}
	return out
}
