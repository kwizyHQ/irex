package diagnostics

import "strings"

// Reporter collects diagnostics from multiple pipeline stages.
// It does NOT print or format anything.
type Reporter struct {
	list []Diagnostic
}

func NewReporter() *Reporter {
	return &Reporter{
		list: make([]Diagnostic, 0),
	}
}

func (r *Reporter) Add(d Diagnostic) {
	r.list = append(r.list, d)
}

func (r *Reporter) Extend(diags []Diagnostic) {
	if len(diags) == 0 {
		return
	}
	r.list = append(r.list, diags...)
}

func (r *Reporter) All() Diagnostics {
	return r.list
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
	if !r.HasErrors() {
		return nil
	}
	return diagnosticsError{diags: r.list}
}

type diagnosticsError struct {
	diags []Diagnostic
}

func (e diagnosticsError) Error() string {
	if len(e.diags) == 0 {
		return "diagnostics error"
	}
	parts := make([]string, 0, len(e.diags))
	for _, d := range e.diags {
		parts = append(parts, d.Message)
	}
	return strings.Join(parts, "; ")
}

func (r *Reporter) HasErrors() bool {
	for _, d := range r.list {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}

func (r *Reporter) HasWarnings() bool {
	for _, d := range r.list {
		if d.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

func (r *Reporter) Error(message string, Range Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    Range,
		Severity: SeverityError,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}

func (r *Reporter) Warn(message string, Range Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    Range,
		Severity: SeverityWarning,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}

func (r *Reporter) Info(message string, Range Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    Range,
		Severity: SeverityInformation,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}

func (r *Reporter) Hint(message string, Range Range, code string, source string) {
	r.Add(Diagnostic{
		Range:    Range,
		Severity: SeverityHint,
		Message:  message,
		Source:   source,
		Code:     code,
	})
}
