// Package tags provides port tagging — associating human-readable labels
// with specific ports or port ranges for richer alert output.
package tags

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Tag associates a label with a set of ports or a port range.
type Tag struct {
	Label string `json:"label"`
	Ports []int  `json:"ports,omitempty"`
	From  int    `json:"from,omitempty"`
	To    int    `json:"to,omitempty"`
}

// Map holds all configured tags keyed by label.
type Map struct {
	tags []Tag
}

// New returns an empty Map.
func New() *Map { return &Map{} }

// Load reads tags from a JSON file.
func Load(path string) (*Map, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, fmt.Errorf("tags: read %s: %w", path, err)
	}
	var tags []Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("tags: parse %s: %w", path, err)
	}
	return &Map{tags: tags}, nil
}

// Lookup returns all labels that match the given port.
func (m *Map) Lookup(port int) []string {
	var labels []string
	for _, t := range m.tags {
		if matchesTag(t, port) {
			labels = append(labels, t.Label)
		}
	}
	sort.Strings(labels)
	return labels
}

// Add appends a tag to the map.
(t Tag) {
	m.tags = append(m.tags, t)
}

// Len returns the number of configured tags.
func (m *Map) Len() int { return len(m.tags) }

func matchesTag(t Tag, port int) bool {
	for _, p := range t.Ports {
		if p == port {
			return true
		}
	}
	if t.From > 0 && t.To > 0 {
		return port >= t.From && port <= t.To
	}
	return false
}
