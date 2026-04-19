// Package portrank assigns a risk/priority rank to ports based on
// well-known service classifications and user-defined overrides.
package portrank

import "sort"

// Rank represents the priority level of a port.
type Rank int

const (
	RankLow    Rank = 1
	RankMedium Rank = 2
	RankHigh   Rank = 3
	RankCritical Rank = 4
)

var rankNames = map[Rank]string{
	RankLow:      "low",
	RankMedium:   "medium",
	RankHigh:     "high",
	RankCritical: "critical",
}

func (r Rank) String() string {
	if s, ok := rankNames[r]; ok {
		return s
	}
	return "unknown"
}

// builtin maps well-known ports to ranks.
var builtin = map[int]Rank{
	22:   RankCritical, // SSH
	23:   RankCritical, // Telnet
	3389: RankCritical, // RDP
	80:   RankHigh,
	443:  RankHigh,
	8080: RankHigh,
	8443: RankHigh,
	3306: RankHigh, // MySQL
	5432: RankHigh, // PostgreSQL
	6379: RankMedium, // Redis
	27017: RankMedium, // MongoDB
	53:   RankMedium, // DNS
}

// Ranker ranks ports.
type Ranker struct {
	overrides map[int]Rank
}

// New returns a Ranker with optional overrides.
func New(overrides map[int]Rank) *Ranker {
	if overrides == nil {
		overrides = make(map[int]Rank)
	}
	return &Ranker{overrides: overrides}
}

// Rank returns the rank for a port.
func (r *Ranker) Rank(port int) Rank {
	if rank, ok := r.overrides[port]; ok {
		return rank
	}
	if rank, ok := builtin[port]; ok {
		return rank
	}
	return RankLow
}

// SortByRank sorts ports descending by rank (highest first).
func (r *Ranker) SortByRank(ports []int) []int {
	out := make([]int, len(ports))
	copy(out, ports)
	sort.SliceStable(out, func(i, j int) bool {
		return r.Rank(out[i]) > r.Rank(out[j])
	})
	return out
}

// FilterByMinRank returns only ports at or above the given rank.
func (r *Ranker) FilterByMinRank(ports []int, min Rank) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if r.Rank(p) >= min {
			out = append(out, p)
		}
	}
	return out
}
