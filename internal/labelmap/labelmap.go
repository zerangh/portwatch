// Package labelmap maps port numbers to human-readable service labels.
package labelmap

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LabelMap holds port-to-label mappings.
type LabelMap struct {
	labels map[int]string
}

// New an empty LabelMap.
func New() *LabelMap {
	return &LabelMap{labels: make(map[int]string)}
}
 a label file where each line is "port label", e.g. "80 http".
// Lines beginning with '#' are treated as comments.
func Load(path string) (*LabelMap, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, fmt.Errorf("labelmap: open %s: %w", path, err)
	}
	defer f.Close()

	lm := New()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		lm.labels[port] = strings.TrimSpace(parts[1])
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("labelmap: read %s: %w", path, err)
	}
	return lm, nil
}

// Lookup returns the label for a port, or an empty string if not found.
func (lm *LabelMap) Lookup(port int) string {
	return lm.labels[port]
}

// Set adds or overwrites a label for a port.
func (lm *LabelMap) Set(port int, label string) {
	lm.labels[port] = label
}

// Annotate returns a map of port -> label for all given ports that have a label.
func (lm *LabelMap) Annotate(ports []int) map[int]string {
	out := make(map[int]string)
	for _, p := range ports {
		if l := lm.labels[p]; l != "" {
			out[p] = l
		}
	}
	return out
}
