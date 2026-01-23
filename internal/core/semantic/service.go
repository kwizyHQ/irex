package semantic

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// ValidateServiceAST checks that all model names referenced in serviceAst are defined in schemaAst.
// Returns a slice of diagnostics for any missing models.
func CheckServiceSemantic(serviceAst *symbols.ServiceDefinition, schemaAst *symbols.ModelsSpec) []diagnostics.Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}
	source := "validate.service.modelref"

	// Build a set of all model names in schemaAst
	modelNames := map[string]struct{}{}
	if schemaAst != nil && schemaAst.ModelsBlock != nil {
		for _, m := range schemaAst.ModelsBlock.Models {
			modelNames[m.Name] = struct{}{}
		}
	}

	// Helper to check a Service and its nested services recursively
	var checkService func(s symbols.Service)
	checkService = func(s symbols.Service) {
		if s.Model != "" {
			if _, ok := modelNames[s.Model]; !ok {
				reporter.Error("Service '"+s.Name+"' references undefined model '"+s.Model+"'", zeroRange,
					"service.model.not_found", source)
			}
		}
		// Recurse into nested services
		for _, nested := range s.Services {
			checkService(nested)
		}
	}

	// Check all top-level services
	if serviceAst != nil && serviceAst.Services != nil {
		for _, svc := range serviceAst.Services.Services {
			checkService(svc)
		}
	}

	return reporter.All()
}
