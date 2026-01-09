package pipeline

import (
	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/core/validate"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

type TemplateOptions struct {
	Path string
}

type TemplateInfo = symbols.TemplateBlock

type TemplateRegistry struct {
	Templates []TemplateInfo
}

func BuildTemplate(opts TemplateOptions) (*TemplateRegistry, error) {
	r := diagnostics.NewReporter()

	// ---------------- Parse Template AST ----------------
	// scan the opts
	var templateDef symbols.TemplateDefinition

	r.Extend(
		ast.ParseHCL(opts.Path, &templateDef).(diagnostics.Diagnostics),
	)

	// validate the template AST

	r.Extend(
		validate.ValidateTemplates(&templateDef),
	)

	if r.HasErrors() {
		return nil, r.Err()
	}

	// map to registry
	Templates := templateDef.Templates

	return &TemplateRegistry{Templates: Templates}, nil

}
