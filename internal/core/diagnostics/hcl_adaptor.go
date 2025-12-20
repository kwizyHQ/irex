package diagnostics

import "github.com/hashicorp/hcl/v2"

// FromHCL converts hcl.Diagnostics into IREX diagnostics.
// This is used ONLY at AST decode time.
func FromHCL(
	errs error,
) []Diagnostic {
	var diags hcl.Diagnostics
	// check if errs is nil (hence no diagnostics)
	if errs == nil {
		return nil
	}

	// now check if errs is iterable then convert to array
	if _, ok := errs.(hcl.Diagnostics); !ok {
		// not iterable, return single diagnostic
		diags.Append(errs.(*hcl.Diagnostic))
	} else {
		diags = errs.(hcl.Diagnostics)
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

		out = append(out, Diagnostic{
			Severity: sev,
			Message:  d.Summary,
			Range: Range{
				Start: Position{
					Line:   d.Subject.Start.Line,
					Column: d.Subject.Start.Column,
					Byte:   d.Subject.Start.Byte,
				},
				End: Position{
					Line:   d.Subject.End.Line,
					Column: d.Subject.End.Column,
					Byte:   d.Subject.End.Byte,
				},
			},
			Source: d.Subject.Filename,
		})
	}

	return out
}

// use a reporter to convert and collect diagnostics
func (r *Reporter) FromHCL(
	diags error,
) {
	r.Extend(FromHCL(diags))
}
