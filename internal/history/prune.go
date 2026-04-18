package history

import (
	"os"
	"encoding/json"
)

// Prune applies the given retention policy to the history file at path,
// removing entries that exceed age or count limits.
func Prune(path string, policy RetentionPolicy) (removed int, err error) {
	h, err := New(path)
	if err != nil {
		return 0, err
	}

	entries, err := h.Load()
	if err != nil {
		return 0, err
	}

	pruned := policy.Apply(entries)
	removed = len(entries) - len(pruned)

	if removed == 0 {
		return 0, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return removed, json.NewEncoder(f).Encode(pruned)
}
