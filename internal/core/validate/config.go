package validate

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// ValidateConfig performs semantic checks on a ConfigDefinition struct and returns diagnostics for all issues found.
func ValidateConfig(cfg *symbols.ConfigDefinition) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}

	if cfg.Project == nil {
		reporter.Error("Missing required 'project' block.", zeroRange, "irex.input.required", "project")
		return reporter.All()
	}
	p := cfg.Project

	if p.Name == "" {
		reporter.Error("Project 'name' is required.", zeroRange, "irex.input.required", "project.name")
	}
	if p.Version == "" {
		reporter.Error("Project 'version' is required.", zeroRange, "irex.input.required", "project.version")
	}
	if p.Author == "" {
		reporter.Warn("Project 'author' is required.", zeroRange, "irex.input.required", "project.author")
	}
	if p.License == "" {
		reporter.Warn("Project 'license' is required.", zeroRange, "irex.input.required", "project.license")
	}

	if p.Paths == nil {
		reporter.Error("Missing required 'paths' block.", zeroRange, "irex.input.required", "project.paths")
	} else {
		if p.Paths.Specifications == "" {
			reporter.Error("'paths.specifications' is required.", zeroRange, "irex.input.required", "project.paths.specifications")
		}
		if p.Paths.Templates == "" {
			reporter.Warn("'paths.templates' is required.", zeroRange, "irex.input.required", "project.paths.templates")
		}
		if p.Paths.Output == "" {
			reporter.Error("'paths.output' is required.", zeroRange, "irex.input.required", "project.paths.output")
		}
	}

	if p.Generator == nil {
		reporter.Error("Missing required 'generator' block.", zeroRange, "irex.input.required", "project.generator")
	}

	if p.Runtime == nil {
		reporter.Error("Missing required 'runtime' block.", zeroRange, "irex.input.required", "project.runtime")
	} else {
		if p.Runtime.Name == "" {
			reporter.Error("'runtime.name' is required.", zeroRange, "irex.input.required", "project.runtime.name")
		}
		if p.Runtime.Version == "" {
			reporter.Warn("'runtime.version' is required.", zeroRange, "irex.input.required", "project.runtime.version")
		}
		if p.Runtime.Options == nil {
			reporter.Error("Missing required 'runtime.options' block.", zeroRange, "irex.input.required", "project.runtime.options")
		} else {
			if p.Runtime.Options.PackageManager == "" {
				reporter.Warn("'runtime.options.package_manager' is required.", zeroRange, "irex.input.required", "project.runtime.options.package_manager")
			}
			if p.Runtime.Options.Entry == "" {
				reporter.Warn("'runtime.options.entry' is required.", zeroRange, "irex.input.required", "project.runtime.options.entry")
			}
		}
		if p.Runtime.Schema == nil {
			reporter.Error("Missing required 'runtime.schema' block.", zeroRange, "irex.input.required", "project.runtime.schema")
		} else {
			if p.Runtime.Schema.Framework == "" {
				reporter.Error("'runtime.schema.framework' is required.", zeroRange, "irex.input.required", "project.runtime.schema.framework")
			}
			if p.Runtime.Schema.Options == nil {
				reporter.Error("Missing required 'runtime.schema.options' block.", zeroRange, "irex.input.required", "project.runtime.schema.options")
			} else {
				// if p.Runtime.Schema.Options.URI == "" {
				// 	reporter.Warn("'runtime.schema.options.uri' is required.", zeroRange, "irex.input.required", "project.runtime.schema.options.uri")
				// }
				// if p.Runtime.Schema.Options.DB == "" {
				// 	reporter.Warn("'runtime.schema.options.db' is required.", zeroRange, "irex.input.required", "project.runtime.schema.options.db")
				// }
			}
		}
		if p.Runtime.Service == nil {
			reporter.Error("Missing required 'runtime.service' block.", zeroRange, "irex.input.required", "project.runtime.service")
		} else {
			if p.Runtime.Service.Framework == "" {
				reporter.Error("'runtime.service.framework' is required.", zeroRange, "irex.input.required", "project.runtime.service.framework")
			}
			if p.Runtime.Service.Options == nil {
				reporter.Error("Missing required 'runtime.service.options' block.", zeroRange, "irex.input.required", "project.runtime.service.options")
			} else {
				if p.Runtime.Service.Options.Port == 0 {
					reporter.Warn("'runtime.service.options.port' is required and must be > 0.", zeroRange, "irex.input.required", "project.runtime.service.options.port")
				}
				if p.Runtime.Service.Options.Host == "" {
					reporter.Warn("'runtime.service.options.host' is required.", zeroRange, "irex.input.required", "project.runtime.service.options.host")
				}
			}
		}
	}

	if p.Meta == nil {
		reporter.Info("Missing optional 'meta' block.", zeroRange, "irex.input.recommended", "project.meta")
	} else {
		if p.Meta.CreatedAt == "" {
			reporter.Info("'meta.created_at' is recommended.", zeroRange, "irex.input.recommended", "project.meta.created_at")
		}
		if p.Meta.GeneratorVersion == "" {
			reporter.Info("'meta.generator_version' is recommended.", zeroRange, "irex.input.recommended", "project.meta.generator_version")
		}
	}

	return reporter.All()
}
