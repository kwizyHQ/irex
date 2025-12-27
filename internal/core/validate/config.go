package validate

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// ValidateConfig performs semantic checks on a ConfigDefinition struct and returns diagnostics for all issues found.
func ValidateConfig(cfg *symbols.ConfigDefinition) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}
	source := "config"

	if cfg.Project == nil {
		reporter.Error("Missing required 'project' block.", zeroRange, "project.missing", source)
		return reporter.All()
	}
	p := cfg.Project

	if p.Name == "" {
		reporter.Error("Project 'name' is required.", zeroRange, "project.name.required", source)
	}
	if p.Version == "" {
		reporter.Error("Project 'version' is required.", zeroRange, "project.version.required", source)
	}
	if p.Author == "" {
		reporter.Warn("Project 'author' is required.", zeroRange, "project.author.required", source)
	}
	if p.License == "" {
		reporter.Warn("Project 'license' is required.", zeroRange, "project.license.required", source)
	}

	if p.Paths == nil {
		reporter.Error("Missing required 'paths' block.", zeroRange, "project.paths.missing", source)
	} else {
		if p.Paths.Specifications == "" {
			reporter.Error("'paths.specifications' is required.", zeroRange, "project.paths.specifications.required", source)
		}
		if p.Paths.Templates == "" {
			reporter.Warn("'paths.templates' is required.", zeroRange, "project.paths.templates.required", source)
		}
		if p.Paths.Output == "" {
			reporter.Error("'paths.output' is required.", zeroRange, "project.paths.output.required", source)
		}
	}

	if p.Generator == nil {
		reporter.Error("Missing required 'generator' block.", zeroRange, "project.generator.missing", source)
	}

	if p.Runtime == nil {
		reporter.Error("Missing required 'runtime' block.", zeroRange, "project.runtime.missing", source)
	} else {
		if p.Runtime.Name == "" {
			reporter.Error("'runtime.name' is required.", zeroRange, "project.runtime.name.required", source)
		}
		if p.Runtime.Version == "" {
			reporter.Warn("'runtime.version' is required.", zeroRange, "project.runtime.version.required", source)
		}
		if p.Runtime.Options == nil {
			reporter.Error("Missing required 'runtime.options' block.", zeroRange, "project.runtime.options.missing", source)
		} else {
			if p.Runtime.Options.PackageManager == "" {
				reporter.Warn("'runtime.options.package_manager' is required.", zeroRange, "project.runtime.options.package_manager.required", source)
			}
			if p.Runtime.Options.Entry == "" {
				reporter.Warn("'runtime.options.entry' is required.", zeroRange, "project.runtime.options.entry.required", source)
			}
		}
		if p.Runtime.Schema == nil {
			reporter.Error("Missing required 'runtime.schema' block.", zeroRange, "project.runtime.schema.missing", source)
		} else {
			if p.Runtime.Schema.Framework == "" {
				reporter.Error("'runtime.schema.framework' is required.", zeroRange, "project.runtime.schema.framework.required", source)
			}
			if p.Runtime.Schema.Options == nil {
				reporter.Error("Missing required 'runtime.schema.options' block.", zeroRange, "project.runtime.schema.options.missing", source)
			} else {
				// if p.Runtime.Schema.Options.URI == "" {
				// 	reporter.Warn("'runtime.schema.options.uri' is required.", zeroRange, "project.runtime.schema.options.uri.required", source)
				// }
				// if p.Runtime.Schema.Options.DB == "" {
				// 	reporter.Warn("'runtime.schema.options.db' is required.", zeroRange, "project.runtime.schema.options.db.required", source)
				// }
			}
		}
		if p.Runtime.Service == nil {
			reporter.Error("Missing required 'runtime.service' block.", zeroRange, "project.runtime.service.missing", source)
		} else {
			if p.Runtime.Service.Framework == "" {
				reporter.Error("'runtime.service.framework' is required.", zeroRange, "project.runtime.service.framework.required", source)
			}
			if p.Runtime.Service.Options == nil {
				reporter.Error("Missing required 'runtime.service.options' block.", zeroRange, "project.runtime.service.options.missing", source)
			} else {
				if p.Runtime.Service.Options.Port == 0 {
					reporter.Warn("'runtime.service.options.port' is required and must be > 0.", zeroRange, "project.runtime.service.options.port.required", source)
				}
				if p.Runtime.Service.Options.Host == "" {
					reporter.Warn("'runtime.service.options.host' is required.", zeroRange, "project.runtime.service.options.host.required", source)
				}
			}
		}
	}

	if p.Meta == nil {
		reporter.Info("Missing optional 'meta' block.", zeroRange, "project.meta.missing", source)
	} else {
		if p.Meta.CreatedAt == "" {
			reporter.Info("'meta.created_at' is recommended.", zeroRange, "project.meta.created_at.recommended", source)
		}
		if p.Meta.GeneratorVersion == "" {
			reporter.Info("'meta.generator_version' is recommended.", zeroRange, "project.meta.generator_version.recommended", source)
		}
	}

	return reporter.All()
}
