package pipeline

import (
	"os"

	"github.com/kwizyHQ/irex/internal/core/ast"
	. "github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/core/validate"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

type TemplateOptions struct {
	Path string
}

type TemplateInfo = TemplateBlock

type TemplateRegistry struct {
	Templates []TemplateInfo
}

func BuildTemplate(opts TemplateOptions) (TemplateRegistry, error) {
	r := diagnostics.NewReporter()

	// ---------------- Parse Template AST ----------------
	// scan the opts
	var templateDef TemplateDefinition
	var registry TemplateRegistry
	r.Extend(
		ast.ParseHCL(opts.Path, &templateDef).(diagnostics.Diagnostics),
	)

	// validate the template AST

	r.Extend(
		validate.ValidateTemplates(&templateDef),
	)

	if r.HasErrors() {
		os.Exit(1)
	}

	// map to registry
	registry.Templates = templateDef.Templates

	return registry, r.All()
}
