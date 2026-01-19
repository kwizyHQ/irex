package diagnostics

import (
	"strings"
	"sync"
)

// Reporter collects diagnostics from multiple pipeline stages.
// It does NOT print or format anything. Reporter is safe for concurrent
// use.
type Reporter struct {
	mu   sync.Mutex
	list []Diagnostic
}

func NewReporter() *Reporter {
	return &Reporter{
		list: make([]Diagnostic, 0),
	}
}

func (r *Reporter) Add(d Diagnostic) {
	r.mu.Lock()
	r.list = append(r.list, d)
	r.mu.Unlock()
}

func (r *Reporter) Extend(diags []Diagnostic) {
	if len(diags) == 0 {
		return
	}
	r.mu.Lock()
	r.list = append(r.list, diags...)
	r.mu.Unlock()
}

func (r *Reporter) All() Diagnostics {
	r.mu.Lock()
	defer r.mu.Unlock()
	// return a copy to avoid races
	out := make([]Diagnostic, len(r.list))
	copy(out, r.list)
	return out
}

func (dc Diagnostics) Error() string {
	if len(dc) == 0 {
		return "no diagnostics"
	}
	parts := make([]string, 0, len(dc))
	for _, d := range dc {
		parts = append(parts, d.Message)
	}
	return strings.Join(parts, "; ")
}

// Err returns an error representing diagnostics if there are any errors present.
// Returns nil when no diagnostics with SeverityError are recorded.
func (r *Reporter) Err() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.hasErrorsLocked() {
		return nil
	}
	return Diagnostics(r.list)
}

func (r *Reporter) hasErrorsLocked() bool {
	for _, d := range r.list {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}

func (r *Reporter) HasErrors() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.hasErrorsLocked()
}

func (r *Reporter) HasWarnings() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.list {
		if d.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

func (r *Reporter) Error(message string, rng Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityError,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}

func (r *Reporter) Warn(message string, rng Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityWarning,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}

func (r *Reporter) Info(message string, rng Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityInformation,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}

func (r *Reporter) Hint(message string, rng Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityHint,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}
