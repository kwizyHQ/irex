package semantic

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// CheckServiceSemantics performs semantic checks on a ServiceDefinition and returns diagnostics for all issues found.
func CheckServiceSemantics(def *symbols.ServiceDefinition) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}
	source := "service"

	if def == nil {
		reporter.Error("Missing service definition root block.", zeroRange, "service.root.missing", source)
		return reporter.All()
	}

	// --- POLICIES ---
	if def.Policies == nil {
		reporter.Error("Missing required 'policies' block.", zeroRange, "service.policies.missing", source)
	} else {
		presetNames := map[string]struct{}{}
		for _, p := range def.Policies.Presets {
			if p.Name == "" {
				reporter.Error("Policy preset missing name.", zeroRange, "service.policies.preset.name.required", source)
			} else {
				if _, exists := presetNames[p.Name]; exists {
					reporter.Error("Duplicate policy preset name: "+p.Name, zeroRange, "service.policies.preset.name.duplicate", source)
				} else {
					presetNames[p.Name] = struct{}{}
				}
			}
			if p.Scope == "" {
				reporter.Warn("Policy preset '"+p.Name+"' missing scope.", zeroRange, "service.policies.preset.scope.recommended", source)
			}
		}
		for _, c := range def.Policies.Customs {
			if c.Name == "" {
				reporter.Error("Custom policy missing name.", zeroRange, "service.policies.custom.name.required", source)
			}
		}
		for _, g := range def.Policies.Groups {
			if g.Name == "" {
				reporter.Error("Policy group missing name.", zeroRange, "service.policies.group.name.required", source)
			}
			if g.Scope == "" {
				reporter.Error("Policy group '"+g.Name+"' missing scope.", zeroRange, "service.policies.group.scope.required", source)
			}
			if len(g.Policies) == 0 {
				reporter.Warn("Policy group '"+g.Name+"' has no policies.", zeroRange, "service.policies.group.policies.empty", source)
			}
		}
	}

	// --- RATE LIMITS ---
	if def.RateLimits == nil {
		reporter.Error("Missing required 'rate_limits' block.", zeroRange, "service.rate_limits.missing", source)
	} else {
		presetNames := map[string]struct{}{}
		for _, p := range def.RateLimits.Presets {
			if p.Name == "" {
				reporter.Error("Rate limit preset missing name.", zeroRange, "service.rate_limits.preset.name.required", source)
			} else {
				if _, exists := presetNames[p.Name]; exists {
					reporter.Error("Duplicate rate limit preset name: "+p.Name, zeroRange, "service.rate_limits.preset.name.duplicate", source)
				} else {
					presetNames[p.Name] = struct{}{}
				}
			}
			if p.Limit == "" && p.Type != "token_bucket" {
				reporter.Warn("Rate limit preset '"+p.Name+"' missing limit.", zeroRange, "service.rate_limits.preset.limit.recommended", source)
			}
		}
		for _, c := range def.RateLimits.Customs {
			if c.Name == "" {
				reporter.Error("Custom rate limit missing name.", zeroRange, "service.rate_limits.custom.name.required", source)
			}
		}
	}

	// --- SERVICES ---
	if def.Services == nil {
		reporter.Error("Missing required 'services' block.", zeroRange, "service.services.missing", source)
	} else {
		if def.Services.BasePath == "" {
			reporter.Warn("Global 'base_path' is recommended.", zeroRange, "service.services.base_path.recommended", source)
		}
		serviceNames := map[string]struct{}{}
		for _, svc := range def.Services.Services {
			checkServiceBlockSemantics(svc, reporter, source, zeroRange, serviceNames)
		}
		for _, op := range def.Services.Operations {
			if op.Name == "" {
				reporter.Error("Global operation missing name.", zeroRange, "service.services.operation.name.required", source)
			}
			if op.Method == "" {
				reporter.Warn("Operation '"+op.Name+"' missing method.", zeroRange, "service.services.operation.method.recommended", source)
			}
			if op.Path == "" {
				reporter.Warn("Operation '"+op.Name+"' missing path.", zeroRange, "service.services.operation.path.recommended", source)
			}
		}
	}

	return reporter.All()
}

func checkServiceBlockSemantics(svc symbols.Service, reporter *diagnostics.Reporter, source string, rng diagnostics.Range, serviceNames map[string]struct{}) {
	if svc.Name == "" {
		reporter.Error("Service block missing name.", rng, "service.services.service.name.required", source)
		return
	}
	if _, exists := serviceNames[svc.Name]; exists {
		reporter.Error("Duplicate service name: "+svc.Name, rng, "service.services.service.name.duplicate", source)
	} else {
		serviceNames[svc.Name] = struct{}{}
	}
	if svc.Model == "" {
		reporter.Warn("Service '"+svc.Name+"' missing model.", rng, "service.services.service.model.recommended", source)
	}
	if svc.Path == "" {
		reporter.Warn("Service '"+svc.Name+"' missing path.", rng, "service.services.service.path.recommended", source)
	}
	for _, op := range svc.Operations {
		if op.Name == "" {
			reporter.Error("Operation in service '"+svc.Name+"' missing name.", rng, "service.services.service.operation.name.required", source)
		}
		if op.Method == "" {
			reporter.Warn("Operation '"+op.Name+"' in service '"+svc.Name+"' missing method.", rng, "service.services.service.operation.method.recommended", source)
		}
		if op.Path == "" {
			reporter.Warn("Operation '"+op.Name+"' in service '"+svc.Name+"' missing path.", rng, "service.services.service.operation.path.recommended", source)
		}
	}
	for _, child := range svc.Services {
		checkServiceBlockSemantics(child, reporter, source, rng, serviceNames)
	}
}
