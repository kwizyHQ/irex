package pipeline

import (
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/core/diagnostics"
	"github.com/kwizyHQ/irex/internal/core/semantic"
	"github.com/kwizyHQ/irex/internal/core/validate"
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
	schemaPath := filepath.Join(ctx.ConfigAST.Project.Paths.Specifications, "schema")
	// fileName (with extension) inside the schemaPath
	if !checkFileExists(schemaPath) {
		r.Error("Schema path does not exist: "+schemaPath, diagnostics.Range{}, "schema.path_not_found", "pipeline")
		return ctx, r.All()
	}
	var schemaFiles []string
	err = filepath.Walk(schemaPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(path) == ".hcl") {
			schemaFiles = append(schemaFiles, path)
		}
		return nil
	})
	if err != nil {
		r.Error("Error reading schema files: "+err.Error(), diagnostics.Range{}, "schema.read_error", "pipeline")
		return ctx, r.All()
	}
	// now parse all schema files and build combined ModelsAST
	var schemaContainsError bool
	for _, filePath := range schemaFiles {
		parsedAST, err := ast.ParseHCLCommon(filePath, "schema")
		if err != nil {
			r.FromHCL(err)
			schemaContainsError = true
			continue
		}
		schemaAST := parsedAST.(*ast.SchemaAST)
		// if ctx.SchemaAST is nil, initialize it
		if ctx.SchemaAST == nil {
			ctx.SchemaAST = schemaAST
		} else {
			// merge ModelsBlock
			ctx.SchemaAST.ModelsBlock.Models = append(ctx.SchemaAST.ModelsBlock.Models, schemaAST.ModelsBlock.Models...)
		}
	}
	if schemaContainsError {
		return ctx, r.All()
	}

	servicesPath := filepath.Join(ctx.ConfigAST.Project.Paths.Specifications, "service")
	if !checkFileExists(servicesPath) {
		r.Error("Services file does not exist.", diagnostics.Range{}, "project.paths.specs", "pipeline")
		return ctx, r.All()
	}
	// ToDo: support multiple service files like schema
	// for now, just parse single service file pick first available .hcl file in servicesPath
	var serviceFile string
	err = filepath.Walk(servicesPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(path) == ".hcl") {
			serviceFile = path
			return filepath.SkipDir // stop after first file
		}
		return nil
	})
	if err != nil {
		r.Error("Error reading service files: "+err.Error(), diagnostics.Range{}, "service.read_error", "pipeline")
		return ctx, r.All()
	}
	serviceAst, err := ast.ParseHCLCommon(serviceFile, "service")
	if err != nil {
		r.FromHCL(err)
		return ctx, r.All()
	}
	ctx.ServicesAST = serviceAst.(*ast.ServicesAST)

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
	r.Extend(
		semantic.CheckServiceSemantics(ctx.ServicesAST),
	)
	r.Extend(
		semantic.CheckSchemaSemantics(ctx.SchemaAST),
	)

	if r.HasErrors() {
		return ctx, r.All()
	}

	// ---------------- Cross Validation: Service Model References ----------------
	// Validate that all service model references exist in schema
	r.Extend(validate.ValidateServiceAST(ctx.ServicesAST, ctx.SchemaAST))

	if r.HasErrors() {
		return ctx, r.All()
	}

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
