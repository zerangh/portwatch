// Package filter provides port filtering utilities for portwatch.
// Filters can be chained to exclude or include ports based on rules.
package filter

// Rule defines a port filtering rule.
type Rule struct {
	ExcludePorts []int
	ExcludeRanges []Range
}

// Range represents an inclusive port range.
type Range struct {
	Low  int
	High int
}

// Filter applies rules to a list of ports and returns the filtered result.
type Filter struct {
	rule Rule
}

// New creates a new Filter with the given rule.
func New(rule Rule) *Filter {
	return &Filter{rule: rule}
}

// Apply returns only ports not excluded by the filter rule.
func (f *Filter) Apply(ports []int) []int {
	excluded := make(map[int]struct{}, len(f.rule.ExcludePorts))
	for _, p := range f.rule.ExcludePorts {
		excluded[p] = struct{}{}
	}

	var result []int
	for _, p := range ports {
		if _, ok := excluded[p]; ok {
			continue
		}
		if f.inExcludedRange(p) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func (f *Filter) inExcludedRange(port int) bool {
	for _, r := range f.rule.ExcludeRanges {
		if port >= r.Low && port <= r.High {
			return true
		}
	}
	return false
}
