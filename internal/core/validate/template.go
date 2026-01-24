package validate

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

func ValidateTemplates(def *symbols.TemplateDefinition) []diagnostics.Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}

	if def == nil {
		reporter.Error("Missing template definition root block.", zeroRange, "irex.input.required", "template")
		return reporter.All()
	}

	// check for duplicate template names
	templateNames := map[string]struct{}{}
	for _, t := range def.Templates {
		if t.Name == "" {
			reporter.Error("Template missing name.", zeroRange, "irex.input.required", "template.name")
		} else {
			if _, exists := templateNames[t.Name]; exists {
				reporter.Error("Duplicate template name: "+t.Name, zeroRange, "irex.input.duplicate", "template.name")
			} else {
				templateNames[t.Name] = struct{}{}
			}
		}
	}

	// check for mode validity (valid modes: "single", "per-item")
	for _, t := range def.Templates {
		if t.Mode != "" && t.Mode != "single" && t.Mode != "per-item" {
			reporter.Error("Invalid template mode '"+t.Mode+"' for template '"+t.Name+"'. Valid modes are 'single' and 'per-item'.", zeroRange, "irex.input.invalid", "template.mode")
		}
		// check if data is set
		if t.Data == "" {
			reporter.Error("Template '"+t.Name+"' has no data defined.", zeroRange, "irex.input.recommended", "template.data")
		}
		// check if output is set
		if t.Output == "" {
			reporter.Error("Template '"+t.Name+"' has no output defined.", zeroRange, "irex.input.recommended", "template.output")
		}
	}

	//TODO
	// now check if template name (i.e. template path) is valid and exists

	return reporter.All()
}
