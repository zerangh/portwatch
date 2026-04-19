// Package portgroup groups ports into named categories (e.g. "web", "db").
package portgroup

import "sort"

// Group maps a category name to a set of ports.
type Group struct {
	Name  string
	Ports []int
}

// Registry holds named port groups.
type Registry struct {
	groups map[string][]int
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string][]int)}
}

// Define registers a named group with the given ports.
func (r *Registry) Define(name string, ports []int) {
	cp := make([]int, len(ports))
	copy(cp, ports)
	sort.Ints(cp)
	r.groups[name] = cp
}

// Lookup returns the ports for a group name and whether it exists.
func (r *Registry) Lookup(name string) ([]int, bool) {
	p, ok := r.groups[name]
	return p, ok
}

// Classify returns all group names that contain the given port.
func (r *Registry) Classify(port int) []string {
	var names []string
	for name, ports := range r.groups {
		for _, p := range ports {
			if p == port {
				names = append(names, name)
				break
			}
		}
	}
	sort.Strings(names)
	return names
}

// All returns all defined groups sorted by name.
func (r *Registry) All() []Group {
	var out []Group
	for name, ports := range r.groups {
		out = append(out, Group{Name: name, Ports: ports})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
