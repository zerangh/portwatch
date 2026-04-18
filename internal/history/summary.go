package history

import "time"

// Summary holds aggregated statistics over a set of history entries.
type Summary struct {
	Total       int
	WithChanges int
	FirstSeen   time.Time
	LastSeen    time.Time
	MostAdded   int
	MostRemoved int
}

// Summarize computes a Summary from the given entries.
func Summarize(entries []Entry) Summary {
	if len(entries) == 0 {
		return Summary{}
	}

	s := Summary{
		Total:     len(entries),
		FirstSeen: entries[0].Timestamp,
		LastSeen:  entries[len(entries)-1].Timestamp,
	}

	for _, e := range entries {
		if !e.Diff.IsEmpty() {
			s.WithChanges++
		}
		if a := len(e.Diff.Added); a > s.MostAdded {
			s.MostAdded = a
		}
		if r := len(e.Diff.Removed); r > s.MostRemoved {
			s.MostRemoved = r
		}
	}

	return s
}

// ChangeRate returns the fraction of entries that contained port changes.
func (s Summary) ChangeRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.WithChanges) / float64(s.Total)
}
