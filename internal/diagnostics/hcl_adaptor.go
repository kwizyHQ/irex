package diagnostics

import "github.com/hashicorp/hcl/v2"

// FromHCL converts hcl.Diagnostics into IREX diagnostics.
// This is used ONLY at AST decode time.
func FromHCL(
	errs error,
) []Diagnostic {
	var diags hcl.Diagnostics
	if errs == nil {
		return nil
	}

	switch v := errs.(type) {
	case hcl.Diagnostics:
		diags = v
	case *hcl.Diagnostic:
		diags = hcl.Diagnostics{v}
	default:
		// unknown error type: return a single diagnostic with the error message
		return []Diagnostic{{
			Severity: SeverityError,
			Message:  errs.Error(),
		}}
	}

	if len(diags) == 0 {
		return nil
	}

	out := make([]Diagnostic, 0, len(diags))

	for _, d := range diags {
		sev := SeverityError
		if d.Severity == hcl.DiagWarning {
			sev = SeverityWarning
		}

		diag := Diagnostic{
			Severity: sev,
			Message:  d.Summary,
		}

		// populate Source and location information if available
		if d.Subject != nil {
			diag.Source = d.Subject.Filename
			diag.Range = Range{}
			// guard Subject.Start/End (hcl.Pos has Line/Column/Byte)
			if d.Subject.Start.Line != 0 || d.Subject.Start.Column != 0 || d.Subject.Start.Byte != 0 {
				diag.Range.Start = Position{
					Line:   d.Subject.Start.Line,
					Column: d.Subject.Start.Column,
					Byte:   d.Subject.Start.Byte,
				}
			}
			if d.Subject.End.Line != 0 || d.Subject.End.Column != 0 || d.Subject.End.Byte != 0 {
				diag.Range.End = Position{
					Line:   d.Subject.End.Line,
					Column: d.Subject.End.Column,
					Byte:   d.Subject.End.Byte,
				}
			}
		}

		out = append(out, diag)
	}

	return out
}

// use a reporter to convert and collect diagnostics
func (r *Reporter) FromHCL(
	diags error,
) {
	r.Extend(FromHCL(diags))
}
