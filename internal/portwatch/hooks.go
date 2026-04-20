package portwatch

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/example/portwatch/internal/portdiff"
)

// HookEvent represents the type of lifecycle event fired to hooks.
type HookEvent string

const (
	HookBeforeScan HookEvent = "before_scan"
	HookAfterScan  HookEvent = "after_scan"
	HookOnChange   HookEvent = "on_change"
	HookOnError    HookEvent = "on_error"
)

// HookFunc is a callback invoked during pipeline lifecycle events.
type HookFunc func(event HookEvent, diff *portdiff.Diff, err error)

// Hooks manages ordered lifecycle callbacks for the scan pipeline.
type Hooks struct {
	handlers []HookFunc
	w        io.Writer
}

// NewHooks returns a new Hooks registry. If w is nil, os.Stderr is used for
// internal diagnostic output.
func NewHooks(w io.Writer) *Hooks {
	if w == nil {
		w = os.Stderr
	}
	return &Hooks{w: w}
}

// Register appends a hook function to the registry.
func (h *Hooks) Register(fn HookFunc) {
	if fn != nil {
		h.handlers = append(h.handlers, fn)
	}
}

// Len returns the number of registered hooks.
func (h *Hooks) Len() int { return len(h.handlers) }

// Fire invokes all registered hooks for the given event, recovering from
// panics so a misbehaving hook cannot crash the pipeline.
func (h *Hooks) Fire(event HookEvent, diff *portdiff.Diff, err error) {
	for i, fn := range h.handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintf(h.w, "portwatch/hooks: panic in hook %d for event %s: %v\n",
						i, event, r)
				}
			}()
			fn(event, diff, err)
		}()
	}
}

// EventNames returns a sorted, comma-separated list of valid hook event names.
func EventNames() string {
	return strings.Join([]string{
		string(HookBeforeScan),
		string(HookAfterScan),
		string(HookOnChange),
		string(HookOnError),
	}, ", ")
}
