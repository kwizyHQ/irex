package diagnostics

// Range represents a text range (LSP-compatible)
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
	// option Byte int `json:"byte,omitempty"`
	Byte int `json:"byte,omitempty"`
}

// Diagnostic represents a single diagnostic message
type Diagnostic struct {
	Range    Range    `json:"range"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Source   string   `json:"source,omitempty"`
	Code     string   `json:"code,omitempty"`
}

type Diagnostics []Diagnostic

func (d Diagnostics) HasErrors() bool {
	for _, diag := range d {
		if diag.Severity == SeverityError {
			return true
		}
	}
	return false
}

// New creates a new Diagnostic
func New(rng Range, severity Severity, message, source, code string) Diagnostic {
	return Diagnostic{
		Range:    rng,
		Severity: severity,
		Message:  message,
		Source:   source,
		Code:     code,
	}
}
