package service

import (
	"fmt"
)

// ValidationError represents a single validation error
// with a path to the offending field and a message.
type ValidationError struct {
	Path    string
	Message string
}

// Validator interface for reusable validation (LSP and CLI)
type Validator interface {
	Validate(def *ServiceDefinition) []ValidationError
}

// ServiceValidator implements Validator for ServiceDefinition
// Can be reused for LSP and CLI validation
// Usage: v := &ServiceValidator{}; errs := v.Validate(def)
type ServiceValidator struct{}

func (v *ServiceValidator) Validate(def *ServiceDefinition) []ValidationError {
	var errs []ValidationError
	if def.Policies != nil {
		errs = append(errs, validatePolicies(def.Policies)...)
	}
	if def.RateLimits != nil {
		errs = append(errs, validateRateLimits(def.RateLimits)...)
	}
	if def.Services != nil {
		errs = append(errs, validateServices(def.Services, def.Policies, def.RateLimits, "services")...)
	}
	return errs
}

// --- Policy Validation ---
func validatePolicies(p *PoliciesBlock) []ValidationError {
	errs := []ValidationError{}
	presetNames := map[string]struct{}{}
	customNames := map[string]struct{}{}
	for _, preset := range p.Presets {
		presetNames[preset.Name] = struct{}{}
		if preset.Scope != "request" && preset.Scope != "resource" {
			errs = append(errs, ValidationError{
				Path:    fmt.Sprintf("policies.policy[%s].scope", preset.Name),
				Message: "Policy scope must be 'request' or 'resource'",
			})
		}
	}
	for _, custom := range p.Customs {
		customNames[custom.Name] = struct{}{}
		if custom.Scope != "request" && custom.Scope != "resource" {
			errs = append(errs, ValidationError{
				Path:    fmt.Sprintf("policies.custom[%s].scope", custom.Name),
				Message: "Custom policy must define scope as 'request' or 'resource'",
			})
		}
	}
	for _, group := range p.Groups {
		if group.Scope != "request" && group.Scope != "resource" {
			errs = append(errs, ValidationError{
				Path:    fmt.Sprintf("policies.group[%s].scope", group.Name),
				Message: "Policy group scope must be 'request' or 'resource'",
			})
		}
	}
	return errs
}

// --- Rate Limit Validation ---
func validateRateLimits(r *RateLimitsBlock) []ValidationError {
	errs := []ValidationError{}
	presetNames := map[string]struct{}{}
	customNames := map[string]struct{}{}
	for _, preset := range r.Presets {
		presetNames[preset.Name] = struct{}{}
	}
	for _, custom := range r.Customs {
		customNames[custom.Name] = struct{}{}
	}
	return errs
}

// --- Service/Operation Validation ---
func validateServices(s *ServicesBlock, p *PoliciesBlock, r *RateLimitsBlock, path string) []ValidationError {
	errs := []ValidationError{}
	for _, svc := range s.Services {
		svcPath := fmt.Sprintf("%s.service[%s]", path, svc.Name)
		errs = append(errs, validateService(&svc, p, r, svcPath)...)
	}
	for _, op := range s.Operations {
		opPath := fmt.Sprintf("%s.operation[%s]", path, op.Name)
		errs = append(errs, validateApplyBlocks(op.Apply, p, r, opPath)...)
	}
	return errs
}

func validateService(svc *Service, p *PoliciesBlock, r *RateLimitsBlock, path string) []ValidationError {
	errs := []ValidationError{}
	for _, apply := range svc.Apply {
		errs = append(errs, validateApplyBlock(apply, p, r, path)...)
	}
	for _, op := range svc.Operations {
		opPath := fmt.Sprintf("%s.operation[%s]", path, op.Name)
		errs = append(errs, validateApplyBlocks(op.Apply, p, r, opPath)...)
	}
	for _, nested := range svc.Services {
		nestedPath := fmt.Sprintf("%s.service[%s]", path, nested.Name)
		errs = append(errs, validateService(&nested, p, r, nestedPath)...)
	}
	return errs
}

func validateApplyBlocks(applyBlocks []ApplyBlock, p *PoliciesBlock, r *RateLimitsBlock, path string) []ValidationError {
	errs := []ValidationError{}
	for _, apply := range applyBlocks {
		errs = append(errs, validateApplyBlock(apply, p, r, path)...)
	}
	return errs
}

func validateApplyBlock(apply ApplyBlock, p *PoliciesBlock, r *RateLimitsBlock, path string) []ValidationError {
	errs := []ValidationError{}
	switch apply.Type {
	case "policy":
		if !policyExists(apply.Name, p) {
			errs = append(errs, ValidationError{
				Path:    fmt.Sprintf("%s.apply[policy:%s]", path, apply.Name),
				Message: "Referenced policy does not exist",
			})
		}
		// Resource-scoped policies must not apply rate limits
		if isResourceScopedPolicy(apply.Name, p) && len(apply.RateLimits) > 0 {
			errs = append(errs, ValidationError{
				Path:    fmt.Sprintf("%s.apply[policy:%s]", path, apply.Name),
				Message: "Resource-scoped policies must not apply rate limits",
			})
		}
		// All referenced rate limits must exist
		for _, rl := range apply.RateLimits {
			if !rateLimitExists(rl, r) {
				errs = append(errs, ValidationError{
					Path:    fmt.Sprintf("%s.apply[policy:%s].rate_limits[%s]", path, apply.Name, rl),
					Message: "Referenced rate limit does not exist",
				})
			}
		}
	case "rate_limit":
		if !rateLimitExists(apply.Name, r) {
			errs = append(errs, ValidationError{
				Path:    fmt.Sprintf("%s.apply[rate_limit:%s]", path, apply.Name),
				Message: "Referenced rate limit does not exist",
			})
		}
	}
	return errs
}

func policyExists(name string, p *PoliciesBlock) bool {
	for _, preset := range p.Presets {
		if preset.Name == name {
			return true
		}
	}
	for _, custom := range p.Customs {
		if custom.Name == name {
			return true
		}
	}
	for _, group := range p.Groups {
		if group.Name == name {
			return true
		}
	}
	return false
}

func isResourceScopedPolicy(name string, p *PoliciesBlock) bool {
	for _, preset := range p.Presets {
		if preset.Name == name && preset.Scope == "resource" {
			return true
		}
	}
	for _, custom := range p.Customs {
		if custom.Name == name && custom.Scope == "resource" {
			return true
		}
	}
	return false
}

func rateLimitExists(name string, r *RateLimitsBlock) bool {
	for _, preset := range r.Presets {
		if preset.Name == name {
			return true
		}
	}
	for _, custom := range r.Customs {
		if custom.Name == name {
			return true
		}
	}
	return false
}
