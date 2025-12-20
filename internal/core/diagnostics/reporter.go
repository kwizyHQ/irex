package diagnostics

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

func (r *Reporter) All() []Diagnostic {
	return r.list
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
