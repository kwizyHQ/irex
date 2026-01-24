package diagnostics

import (
	"sync"
)

// Reporter collects diagnostics from multiple pipeline stages.
// It does NOT print or format anything. Reporter is safe for concurrent
// use.
type Reporter struct {
	mu       sync.Mutex
	list     []Diagnostic
	filename string
}

func NewReporter() *Reporter {
	return &Reporter{
		list: make([]Diagnostic, 0),
	}
}

func (r *Reporter) SetFilename(fn string) {
	r.mu.Lock()
	r.filename = fn
	r.mu.Unlock()
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

func (r *Reporter) ExtendWithFilename(diags []Diagnostic) {
	if len(diags) == 0 {
		return
	}
	r.mu.Lock()
	for _, d := range diags {
		d.Filename = r.filename
		r.list = append(r.list, d)
	}
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

func (r *Reporter) HasErrors() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.list {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
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

func (r *Reporter) Error(message string, rng Range, code string, hclPath string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityError,
		Message:  message,
		Code:     code,
		HclPath:  hclPath,
		Filename: r.filename,
	})
}

func (r *Reporter) Warn(message string, rng Range, code string, hclPath string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityWarning,
		Message:  message,
		Code:     code,
		HclPath:  hclPath,
		Filename: r.filename,
	})
}

func (r *Reporter) Info(message string, rng Range, code string, hclPath string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityInformation,
		Message:  message,
		Code:     code,
		HclPath:  hclPath,
		Filename: r.filename,
	})
}

func (r *Reporter) Hint(message string, rng Range, code string, hclPath string) {
	r.Add(Diagnostic{
		Range:    rng,
		Severity: SeverityHint,
		Message:  message,
		Code:     code,
		HclPath:  hclPath,
		Filename: r.filename,
	})
}
