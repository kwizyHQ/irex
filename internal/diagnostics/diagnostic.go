package diagnostics

import "fmt"

// Range represents a text range (LSP-compatible)
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
	Byte   int `json:"byte,omitempty"`
}

// Diagnostic represents a single diagnostic message
type Diagnostic struct {
	Range    Range    `json:"range"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Source   string   `json:"source,omitempty"`
	Code     string   `json:"code,omitempty"`
	Filename string   `json:"filename,omitempty"`
	HclPath  string   `json:"hcl_path,omitempty"`
}

type Diagnostics []Diagnostic

func (d *Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s; %s", d.Source, d.Message, d.Code)
}

// error implementation, so that sets of diagnostics can be returned via
// APIs that normally deal in vanilla Go errors.
func (d Diagnostics) Error() string {
	count := len(d)
	switch count {
	case 0:
		return "no diagnostics"
	case 1:
		return d[0].Error()
	default:
		return fmt.Sprintf("%s, and %d other diagnostic(s)", d[0].Error(), count-1)
	}
}
