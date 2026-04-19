package filter

// Chain holds multiple filters applied in sequence.
type Chain struct {
	filters []*Filter
}

// NewChain creates a Chain from the provided filters.
func NewChain(filters ...*Filter) *Chain {
	return &Chain{filters: filters}
}

// Apply passes ports through each filter in order.
func (c *Chain) Apply(ports []int) []int {
	result := ports
	for _, f := range c.filters {
		result = f.Apply(result)
	}
	return result
}

// Len returns the number of filters in the chain.
func (c *Chain) Len() int {
	return len(c.filters)
}
