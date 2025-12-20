package semantic

import (
	"github.com/kwizyHQ/irex/internal/core/diagnostics"
	"github.com/kwizyHQ/irex/internal/core/symbols"
)

type Diagnostic = diagnostics.Diagnostic

// CheckConfigSemantics performs semantic checks on a ConfigDefinition struct and returns diagnostics for all issues found.
func CheckConfigSemantics(cfg *symbols.ConfigDefinition) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{
		Start: diagnostics.Position{Line: 0, Column: 0, Byte: 0},
		End:   diagnostics.Position{Line: 0, Column: 0, Byte: 0},
	}

	if cfg.Project == nil {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityError,
			Message:  "Missing required 'project' block.",
			Source:   "semantic",
			Code:     "project.missing",
		})
		return reporter.All()
	}
	p := cfg.Project

	if p.Name == "" {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityError,
			Message:  "Project 'name' is required.",
			Source:   "semantic",
			Code:     "project.name.required",
		})
	}
	if p.Version == "" {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityError,
			Message:  "Project 'version' is required.",
			Source:   "semantic",
			Code:     "project.version.required",
		})
	}
	if p.Author == "" {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityWarning,
			Message:  "Project 'author' is required.",
			Source:   "semantic",
			Code:     "project.author.required",
		})
	}
	if p.License == "" {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityWarning,
			Message:  "Project 'license' is required.",
			Source:   "semantic",
			Code:     "project.license.required",
		})
	}

	if p.Paths == nil {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityError,
			Message:  "Missing required 'paths' block.",
			Source:   "semantic",
			Code:     "project.paths.missing",
		})
	} else {
		if p.Paths.Specifications == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityError,
				Message:  "'paths.specifications' is required.",
				Source:   "semantic",
				Code:     "project.paths.specifications.required",
			})
		}
		if p.Paths.Templates == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityWarning,
				Message:  "'paths.templates' is required.",
				Source:   "semantic",
				Code:     "project.paths.templates.required",
			})
		}
		if p.Paths.Output == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityError,
				Message:  "'paths.output' is required.",
				Source:   "semantic",
				Code:     "project.paths.output.required",
			})
		}
	}

	if p.Generator == nil {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityError,
			Message:  "Missing required 'generator' block.",
			Source:   "semantic",
			Code:     "project.generator.missing",
		})
	}

	if p.Runtime == nil {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityError,
			Message:  "Missing required 'runtime' block.",
			Source:   "semantic",
			Code:     "project.runtime.missing",
		})
	} else {
		if p.Runtime.Name == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityError,
				Message:  "'runtime.name' is required.",
				Source:   "semantic",
				Code:     "project.runtime.name.required",
			})
		}
		if p.Runtime.Version == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityWarning,
				Message:  "'runtime.version' is required.",
				Source:   "semantic",
				Code:     "project.runtime.version.required",
			})
		}
		if p.Runtime.Options == nil {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityError,
				Message:  "Missing required 'runtime.options' block.",
				Source:   "semantic",
				Code:     "project.runtime.options.missing",
			})
		} else {
			if p.Runtime.Options.PackageManager == "" {
				reporter.Add(Diagnostic{
					Range:    zeroRange,
					Severity: diagnostics.SeverityWarning,
					Message:  "'runtime.options.package_manager' is required.",
					Source:   "semantic",
					Code:     "project.runtime.options.package_manager.required",
				})
			}
			if p.Runtime.Options.Entry == "" {
				reporter.Add(Diagnostic{
					Range:    zeroRange,
					Severity: diagnostics.SeverityWarning,
					Message:  "'runtime.options.entry' is required.",
					Source:   "semantic",
					Code:     "project.runtime.options.entry.required",
				})
			}
		}
		if p.Runtime.Schema == nil {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityError,
				Message:  "Missing required 'runtime.schema' block.",
				Source:   "semantic",
				Code:     "project.runtime.schema.missing",
			})
		} else {
			if p.Runtime.Schema.Framework == "" {
				reporter.Add(Diagnostic{
					Range:    zeroRange,
					Severity: diagnostics.SeverityError,
					Message:  "'runtime.schema.framework' is required.",
					Source:   "semantic",
					Code:     "project.runtime.schema.framework.required",
				})
			}
			if p.Runtime.Schema.Options == nil {
				reporter.Add(Diagnostic{
					Range:    zeroRange,
					Severity: diagnostics.SeverityError,
					Message:  "Missing required 'runtime.schema.options' block.",
					Source:   "semantic",
					Code:     "project.runtime.schema.options.missing",
				})
			} else {
				// if p.Runtime.Schema.Options.URI == "" {
				// 	reporter.Add(Diagnostic{
				// 		Range:    zeroRange,
				// 		Severity: diagnostics.SeverityWarning,
				// 		Message:  "'runtime.schema.options.uri' is required.",
				// 		Source:   "semantic",
				// 		Code:     "project.runtime.schema.options.uri.required",
				// 	})
				// }
				// if p.Runtime.Schema.Options.DB == "" {
				// 	reporter.Add(Diagnostic{
				// 		Range:    zeroRange,
				// 		Severity: diagnostics.SeverityWarning,
				// 		Message:  "'runtime.schema.options.db' is required.",
				// 		Source:   "semantic",
				// 		Code:     "project.runtime.schema.options.db.required",
				// 	})
				// }
			}
		}
		if p.Runtime.Service == nil {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityError,
				Message:  "Missing required 'runtime.service' block.",
				Source:   "semantic",
				Code:     "project.runtime.service.missing",
			})
		} else {
			if p.Runtime.Service.Framework == "" {
				reporter.Add(Diagnostic{
					Range:    zeroRange,
					Severity: diagnostics.SeverityError,
					Message:  "'runtime.service.framework' is required.",
					Source:   "semantic",
					Code:     "project.runtime.service.framework.required",
				})
			}
			if p.Runtime.Service.Options == nil {
				reporter.Add(Diagnostic{
					Range:    zeroRange,
					Severity: diagnostics.SeverityError,
					Message:  "Missing required 'runtime.service.options' block.",
					Source:   "semantic",
					Code:     "project.runtime.service.options.missing",
				})
			} else {
				if p.Runtime.Service.Options.Port == 0 {
					reporter.Add(Diagnostic{
						Range:    zeroRange,
						Severity: diagnostics.SeverityWarning,
						Message:  "'runtime.service.options.port' is required and must be > 0.",
						Source:   "semantic",
						Code:     "project.runtime.service.options.port.required",
					})
				}
				if p.Runtime.Service.Options.Host == "" {
					reporter.Add(Diagnostic{
						Range:    zeroRange,
						Severity: diagnostics.SeverityWarning,
						Message:  "'runtime.service.options.host' is required.",
						Source:   "semantic",
						Code:     "project.runtime.service.options.host.required",
					})
				}
			}
		}
	}

	if p.Meta == nil {
		reporter.Add(Diagnostic{
			Range:    zeroRange,
			Severity: diagnostics.SeverityInformation,
			Message:  "Missing optional 'meta' block.",
			Source:   "semantic",
			Code:     "project.meta.missing",
		})
	} else {
		if p.Meta.CreatedAt == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityInformation,
				Message:  "'meta.created_at' is recommended.",
				Source:   "semantic",
				Code:     "project.meta.created_at.recommended",
			})
		}
		if p.Meta.GeneratorVersion == "" {
			reporter.Add(Diagnostic{
				Range:    zeroRange,
				Severity: diagnostics.SeverityInformation,
				Message:  "'meta.generator_version' is recommended.",
				Source:   "semantic",
				Code:     "project.meta.generator_version.recommended",
			})
		}
	}

	return reporter.All()
}
