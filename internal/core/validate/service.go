package validate

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// ValidateService performs semantic checks on a ServiceDefinition and returns diagnostics for all issues found.
func ValidateService(def *symbols.ServiceDefinition) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}

	if def == nil {
		reporter.Error("Missing service definition root block.", zeroRange, "irex.input.required", "service")
		return reporter.All()
	}

	// --- POLICIES ---
	if def.Policies == nil {
		reporter.Error("Missing required 'policies' block.", zeroRange, "irex.input.required", "service.policies")
	} else {
		presetNames := map[string]struct{}{}
		for _, p := range def.Policies.Presets {
			if p.Name == "" {
				reporter.Error("Policy preset missing name.", zeroRange, "irex.input.required", "service.policies.preset.name")
			} else {
				if _, exists := presetNames[p.Name]; exists {
					reporter.Error("Duplicate policy preset name: "+p.Name, zeroRange, "irex.input.duplicate", "service.policies.preset.name")
				} else {
					presetNames[p.Name] = struct{}{}
				}
			}
			if p.Scope == "" {
				reporter.Warn("Policy preset '"+p.Name+"' missing scope.", zeroRange, "irex.input.recommended", "service.policies.preset.scope")
			}
		}
		for _, c := range def.Policies.Customs {
			if c.Name == "" {
				reporter.Error("Custom policy missing name.", zeroRange, "irex.input.required", "service.policies.custom.name")
			}
		}
		for _, g := range def.Policies.Groups {
			if g.Name == "" {
				reporter.Error("Policy group missing name.", zeroRange, "irex.input.required", "service.policies.group.name")
			}
			if g.Scope == "" {
				reporter.Error("Policy group '"+g.Name+"' missing scope.", zeroRange, "irex.input.required", "service.policies.group.scope")
			}
			if len(g.Policies) == 0 {
				reporter.Warn("Policy group '"+g.Name+"' has no policies.", zeroRange, "irex.input.recommended", "service.policies.group.policies")
			}
		}
	}

	// --- RATE LIMITS ---
	if def.RateLimits == nil {
		reporter.Error("Missing required 'rate_limits' block.", zeroRange, "irex.input.required", "service.rate_limits")
	} else {
		presetNames := map[string]struct{}{}
		for _, p := range def.RateLimits.Presets {
			if p.Name == "" {
				reporter.Error("Rate limit preset missing name.", zeroRange, "irex.input.required", "service.rate_limits.preset.name")
			} else {
				if _, exists := presetNames[p.Name]; exists {
					reporter.Error("Duplicate rate limit preset name: "+p.Name, zeroRange, "irex.input.duplicate", "service.rate_limits.preset.name")
				} else {
					presetNames[p.Name] = struct{}{}
				}
			}
			if p.Limit == "" && p.Type != "token_bucket" {
				reporter.Warn("Rate limit preset '"+p.Name+"' missing limit.", zeroRange, "irex.input.recommended", "service.rate_limits.preset.limit")
			}
		}
		for _, c := range def.RateLimits.Customs {
			if c.Name == "" {
				reporter.Error("Custom rate limit missing name.", zeroRange, "irex.input.required", "service.rate_limits.custom.name")
			}
		}
	}

	// --- SERVICES ---
	if def.Services == nil {
		reporter.Error("Missing required 'services' block.", zeroRange, "irex.input.required", "service.services")
	} else {
		if def.Services.BasePath == "" {
			reporter.Warn("Global 'base_path' is recommended.", zeroRange, "irex.input.recommended", "service.services.base_path")
		}
		serviceNames := map[string]struct{}{}
		for _, svc := range def.Services.Services {
			checkServiceBlockSemantics(svc, reporter, zeroRange, serviceNames)
		}
		for _, op := range def.Services.Operations {
			if op.Name == "" {
				reporter.Error("Global operation missing name.", zeroRange, "irex.input.required", "service.services.operation.name")
			}
			if op.Method == "" {
				reporter.Warn("Operation '"+op.Name+"' missing method.", zeroRange, "irex.input.recommended", "service.services.operation.method")
			}
			if op.Path == "" {
				reporter.Warn("Operation '"+op.Name+"' missing path.", zeroRange, "irex.input.recommended", "service.services.operation.path")
			}
		}
	}

	return reporter.All()
}

func checkServiceBlockSemantics(svc symbols.Service, reporter *diagnostics.Reporter, rng diagnostics.Range, serviceNames map[string]struct{}) {
	if svc.Name == "" {
		reporter.Error("Service block missing name.", rng, "irex.input.required", "service.services.service.name")
		return
	}
	if _, exists := serviceNames[svc.Name]; exists {
		reporter.Error("Duplicate service name: "+svc.Name, rng, "irex.input.duplicate", "service.services.service.name")
	} else {
		serviceNames[svc.Name] = struct{}{}
	}
	if svc.Model == "" {
		reporter.Warn("Service '"+svc.Name+"' missing model.", rng, "irex.input.recommended", "service.services.service.model")
	}
	if svc.Path == "" {
		reporter.Warn("Service '"+svc.Name+"' missing path.", rng, "irex.input.recommended", "service.services.service.path")
	}
	for _, op := range svc.Operations {
		if op.Name == "" {
			reporter.Error("Operation in service '"+svc.Name+"' missing name.", rng, "irex.input.required", "service.services.service.operation.name")
		}
		if op.Method == "" {
			reporter.Warn("Operation '"+op.Name+"' in service '"+svc.Name+"' missing method.", rng, "irex.input.recommended", "service.services.service.operation.method")
		}
		if op.Path == "" {
			reporter.Warn("Operation '"+op.Name+"' in service '"+svc.Name+"' missing path.", rng, "irex.input.recommended", "service.services.service.operation.path")
		}
	}
	for _, child := range svc.Services {
		checkServiceBlockSemantics(child, reporter, rng, serviceNames)
	}
}
