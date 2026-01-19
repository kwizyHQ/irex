package diagnostics

// Range represents a text range (LSP-compatible)
type Range struct {
	Start Position `json:"start,omitempty"`
	End   Position `json:"end,omitempty"`
}

type Position struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
	Byte   int `json:"byte,omitempty"`
}

// Location is a more agnostic location representation that can hold
// filename, byte offsets and optionally LSP-like positions. New code
// should prefer `Location` when a source-agnostic representation is
// needed. `Range` remains for backwards compatibility.
type Location struct {
	Filename  string    `json:"filename,omitempty"`
	StartByte int       `json:"start_byte,omitempty"`
	EndByte   int       `json:"end_byte,omitempty"`
	Start     *Position `json:"start,omitempty"`
	End       *Position `json:"end,omitempty"`
}

// Diagnostic represents a single diagnostic message
type Diagnostic struct {
	Range    Range     `json:"range,omitempty"`
	Location *Location `json:"location,omitempty"`
	Severity Severity  `json:"severity"`
	Message  string    `json:"message"`
	Source   string    `json:"source,omitempty"`
	Code     string    `json:"code,omitempty"`
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

// NewDiagnostic creates a new Diagnostic. New is kept for compatibility.
func NewDiagnostic(rng Range, severity Severity, message, source, code string) Diagnostic {
	return Diagnostic{
		Range:    rng,
		Severity: severity,
		Message:  message,
		Source:   source,
		Code:     code,
	}
}

//go:deprecated: New is kept for compatibility; prefer NewDiagnostic.
func New(rng Range, severity Severity, message, source, code string) Diagnostic {
	return NewDiagnostic(rng, severity, message, source, code)
}
