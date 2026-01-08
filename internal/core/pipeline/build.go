package pipeline

import (
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/assemble"
	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/core/normalize"
	"github.com/kwizyHQ/irex/internal/core/semantic"
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/core/validate"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

type BuildOptions struct {
	ConfigPath string
}

func Build(opts BuildOptions) (*shared.IRBundle, error) {
	r := diagnostics.NewReporter()
	ctx := &shared.BuildContext{
		ConfigAST: &shared.ConfigAST{},
		SchemaAST: &shared.SchemaAST{
			ModelsBlock: &symbols.ModelsBlock{
				Models: make([]symbols.Model, 0),
			},
		},
		ServicesAST: &shared.ServicesAST{},
		IR:          &shared.IRBundle{},
	}

	// ------------------- Config AST Decode ----------------
	r.Extend(ast.ParseHCL(opts.ConfigPath, ctx.ConfigAST).(diagnostics.Diagnostics))

	r.Extend(
		validate.ValidateConfig(ctx.ConfigAST),
	)

	if r.HasErrors() {
		return nil, r.Err()
	}

	// ---------------- Other AST Decode ----------------
	schemaPath := filepath.Join(ctx.ConfigAST.Project.Paths.Specifications, "schema")
	files, _ := filepath.Glob(filepath.Join(schemaPath, "*.hcl"))

	var schemaContainsError bool
	for _, path := range files {
		var spec symbols.ModelsSpec
		// We assume ParseHCL now handles the pointer internally
		if diags := ast.ParseHCL(path, &spec).(diagnostics.Diagnostics); diags.HasErrors() {
			schemaContainsError = true
			r.Extend(diags)
			continue
		}
		// Direct append using the spread operator (...)
		if spec.ModelsBlock != nil {
			ctx.SchemaAST.ModelsBlock.Models = append(ctx.SchemaAST.ModelsBlock.Models, spec.ModelsBlock.Models...)
		}
	}

	if schemaContainsError {
		return nil, r.All()
	}

	servicesPath := filepath.Join(ctx.ConfigAST.Project.Paths.Specifications, "service")
	files, err := filepath.Glob(filepath.Join(servicesPath, "*.hcl"))
	if err != nil || len(files) == 0 {
		r.Error("Warning we couln't found any service files, please add some.", diagnostics.Range{}, "service.read_error", "pipeline")
		return nil, r.Err()
	}
	r.Extend(
		ast.ParseHCL(files[0], ctx.ServicesAST).(diagnostics.Diagnostics),
	)

	// if reporter.HasErrors() {
	// 	return nil, reporter.All()
	// }

	// // ---------------- Registry ----------------

	// reg := registry.New()
	// if err := reg.RegisterModels(ctx.ModelsAST); err != nil {
	// 	reporter.Add(diagnostics.Diagnostic{
	// 		Severity: diagnostics.SeverityError,
	// 		Message:  err.Error(),
	// 		Source:   "registry",
	// 	})
	// }

	// if err := reg.RegisterServices(ctx.ServicesAST); err != nil {
	// 	reporter.Add(diagnostics.Diagnostic{
	// 		Severity: diagnostics.SeverityError,
	// 		Message:  err.Error(),
	// 		Source:   "registry",
	// 	})
	// }

	// ctx.Registry = reg

	if r.HasErrors() {
		return nil, r.All()
	}

	// ---------------- Validations ----------------

	r.Extend(
		validate.ValidateService(ctx.ServicesAST),
	)
	r.Extend(
		validate.ValidateSchema(ctx.SchemaAST),
	)

	if r.HasErrors() {
		return nil, r.All()
	}

	// ---------------- Cross Validation: Semantic checks ----------------
	// Validate that all service model references exist in schema
	r.Extend(semantic.CheckServiceSemantic(ctx.ServicesAST, ctx.SchemaAST))

	if r.HasErrors() {
		return nil, r.All()
	}

	// // ---------------- Normalize ----------------

	normalize.NormalizeServiceAST(ctx.ServicesAST)

	// ---------------- IR Build ----------------
	err = assemble.ProjectIR(ctx)

	if err != nil {
		r.Error("IR Build error: "+err.Error(), diagnostics.Range{}, "ir.build_error", "pipeline")
	}

	return ctx.IR, r.Err()
}
