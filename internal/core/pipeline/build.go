package pipeline

import (
	"os"

	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/core/diagnostics"
	"github.com/kwizyHQ/irex/internal/core/semantic"
)

func Build(opts BuildOptions) (*BuildContext, []diagnostics.Diagnostic) {
	r := diagnostics.NewReporter()
	ctx := &BuildContext{}

	if !checkFileExists(opts.ConfigPath) {
		r.Error("Config file does not exist.", diagnostics.Range{}, "config.not_found", "pipeline")
		return ctx, r.All()
	}

	// ---------------- Config AST Decode ----------------
	configAST, err := ast.ParseHCLCommon(opts.ConfigPath, "config")

	if err != nil {
		r.FromHCL(err)
	}

	ctx.ConfigAST = configAST.(*ast.ConfigAST)
	// ---------------- End Config AST Decode ----------------

	if r.HasErrors() {
		// return early if there are errors in config parsing as other steps depend on it
		return ctx, r.All()
	}

	// ---------------- Other AST Decode ----------------

	// ctx.ModelsAST = &ast.ModelsAST{}
	// reporter.Extend(
	// 	diagnostics.FromHCL(
	// 		ast.DecodeFile(opts.ModelsPath, ctx.ModelsAST),
	// 		"ast",
	// 	),
	// )

	// ctx.ServicesAST = &ast.ServicesAST{}
	// reporter.Extend(
	// 	diagnostics.FromHCL(
	// 		ast.DecodeFile(opts.ServicesPath, ctx.ServicesAST),
	// 		"ast",
	// 	),
	// )

	// if reporter.HasErrors() {
	// 	return ctx, reporter.All()
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
		return ctx, r.All()
	}

	// ---------------- Semantic ----------------

	r.Extend(
		semantic.CheckConfigSemantics(ctx.ConfigAST),
	)
	// reporter.Extend(
	// 	semantic.ValidateServices(ctx.ServicesAST),
	// )

	if r.HasErrors() {
		return ctx, r.All()
	}

	// // ---------------- Cross Validation ----------------

	// reporter.Extend(
	// 	validate.ValidateServiceModelRefs(reg),
	// )

	// if reporter.HasErrors() {
	// 	return ctx, reporter.All()
	// }

	// // ---------------- Normalize ----------------

	// norm := normalize.New(reg)
	// irBundle, normDiags := norm.BuildIR()
	// reporter.Extend(normDiags)

	// ctx.IR = irBundle

	return ctx, r.All()
}

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
