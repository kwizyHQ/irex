package config

// ValidationError represents a single validation error
// with a path to the offending field and a message.
type ValidationError struct {
	Path    string
	Message string
}

// Validator interface for reusable validation (LSP and CLI)
type Validator interface {
	Validate(def *ConfigDefinition) []ValidationError
}

// ConfigValidator implements Validator for ConfigDefinition
// Can be reused for LSP and CLI validation
// Usage: v := &ConfigValidator{}; errs := v.Validate(def)
type ConfigValidator struct{}

func (v *ConfigValidator) Validate(def *ConfigDefinition) []ValidationError {
	var errs []ValidationError
	if def.Project == nil {
		errs = append(errs, ValidationError{
			Path:    "project",
			Message: "Project block must be provided",
		})
		return errs
	}
	p := def.Project
	if p.Name == "" {
		errs = append(errs, ValidationError{
			Path:    "project.name",
			Message: "Project name must be provided",
		})
	}
	if p.Paths == nil {
		errs = append(errs, ValidationError{
			Path:    "project.paths",
			Message: "Paths block must be provided",
		})
	} else {
		if p.Paths.Specifications == "" {
			errs = append(errs, ValidationError{
				Path:    "project.paths.specifications",
				Message: "Path to specifications must be provided",
			})
		}
		if p.Paths.Output == "" {
			errs = append(errs, ValidationError{
				Path:    "project.paths.output",
				Message: "Output path must be provided",
			})
		}
	}
	if p.Generator == nil {
		errs = append(errs, ValidationError{
			Path:    "project.generator",
			Message: "Generator block must be provided",
		})
	}
	if p.Runtime == nil {
		errs = append(errs, ValidationError{
			Path:    "project.runtime",
			Message: "Runtime block must be provided",
		})
	} else {
		r := p.Runtime
		if r.Name == "" {
			errs = append(errs, ValidationError{
				Path:    "project.runtime.name",
				Message: "Runtime name must be provided",
			})
		}
		if r.Options == nil {
			errs = append(errs, ValidationError{
				Path:    "project.runtime.options",
				Message: "Runtime options block must be provided",
			})
		}
		if r.Schema == nil {
			errs = append(errs, ValidationError{
				Path:    "project.runtime.schema",
				Message: "Runtime schema block must be provided",
			})
		} else {
			if r.Schema.Framework == "" {
				errs = append(errs, ValidationError{
					Path:    "project.runtime.schema.framework",
					Message: "Schema framework must be provided",
				})
			}
			if r.Schema.Options == nil {
				errs = append(errs, ValidationError{
					Path:    "project.runtime.schema.options",
					Message: "Schema options block must be provided",
				})
			}
		}
		if r.Service == nil {
			errs = append(errs, ValidationError{
				Path:    "project.runtime.service",
				Message: "Runtime service block must be provided",
			})
		} else {
			if r.Service.Framework == "" {
				errs = append(errs, ValidationError{
					Path:    "project.runtime.service.framework",
					Message: "Service framework must be provided",
				})
			}
			if r.Service.Options == nil {
				errs = append(errs, ValidationError{
					Path:    "project.runtime.service.options",
					Message: "Service options block must be provided",
				})
			}
		}
	}
	if p.Meta == nil {
		errs = append(errs, ValidationError{
			Path:    "project.meta",
			Message: "Meta block must be provided",
		})
	}
	return errs
}
