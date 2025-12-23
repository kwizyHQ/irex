package normalize

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
)

// MergeServiceDefaults merges parent and child ServiceDefaults, with child taking precedence.
func MergeServiceDefaults(parent, child *symbols.ServiceDefaults) *symbols.ServiceDefaults {
	if parent == nil && child == nil {
		return nil
	}
	if parent == nil {
		return child
	}
	if child == nil {
		return parent
	}
	merged := *parent // shallow copy
	// For each field, if child has non-nil/non-zero, use child
	if child.Pagination != nil {
		merged.Pagination = child.Pagination
	}
	if child.Expose != nil {
		merged.Expose = child.Expose
	}
	if child.SoftDelete != nil {
		merged.SoftDelete = child.SoftDelete
	}
	if len(child.CrudOperations) > 0 {
		merged.CrudOperations = child.CrudOperations
	}
	if len(child.BatchOperations) > 0 {
		merged.BatchOperations = child.BatchOperations
	}
	if len(child.Middlewares) > 0 {
		merged.Middlewares = child.Middlewares
	}
	if len(child.Sorting) > 0 {
		merged.Sorting = child.Sorting
	}
	if len(child.Filtering) > 0 {
		merged.Filtering = child.Filtering
	}
	if len(child.Search) > 0 {
		merged.Search = child.Search
	}
	return &merged
}

// MergeRateLimitPresets merges parent and child RateLimitPresets, child overrides by name
func MergeRateLimitPresets(parent, child []symbols.RateLimitPreset) []symbols.RateLimitPreset {
	merged := make([]symbols.RateLimitPreset, 0, len(parent)+len(child))
	byName := map[string]symbols.RateLimitPreset{}
	for _, p := range parent {
		byName[p.Name] = p
	}
	for _, c := range child {
		byName[c.Name] = c // child overrides
	}
	for _, v := range byName {
		merged = append(merged, v)
	}
	return merged
}

// MergePolicyPresets merges parent and child PolicyPresets, child overrides by name
func MergePolicyPresets(parent, child []symbols.PolicyPreset) []symbols.PolicyPreset {
	merged := make([]symbols.PolicyPreset, 0, len(parent)+len(child))
	byName := map[string]symbols.PolicyPreset{}
	for _, p := range parent {
		byName[p.Name] = p
	}
	for _, c := range child {
		byName[c.Name] = c // child overrides
	}
	for _, v := range byName {
		merged = append(merged, v)
	}
	return merged
}

// NormalizeServiceAST recursively merges defaults, rate_limits, and policies for all services.
func NormalizeServiceAST(def *symbols.ServiceDefinition) *symbols.ServiceDefinition {
	if def == nil || def.Services == nil {
		return def
	}

	// Merge policies presets recursively
	var mergePolicies func(parent *symbols.PoliciesBlock, child *symbols.PoliciesBlock) *symbols.PoliciesBlock
	mergePolicies = func(parent, child *symbols.PoliciesBlock) *symbols.PoliciesBlock {
		if parent == nil && child == nil {
			return nil
		}
		if parent == nil {
			return child
		}
		if child == nil {
			return parent
		}
		merged := *parent
		merged.Presets = MergePolicyPresets(parent.Presets, child.Presets)
		merged.Customs = append(parent.Customs, child.Customs...)
		merged.Groups = append(parent.Groups, child.Groups...)
		return &merged
	}

	// Merge rate limit presets recursively
	var mergeRateLimits func(parent *symbols.RateLimitsBlock, child *symbols.RateLimitsBlock) *symbols.RateLimitsBlock
	mergeRateLimits = func(parent, child *symbols.RateLimitsBlock) *symbols.RateLimitsBlock {
		if parent == nil && child == nil {
			return nil
		}
		if parent == nil {
			return child
		}
		if child == nil {
			return parent
		}
		merged := *parent
		merged.Presets = MergeRateLimitPresets(parent.Presets, child.Presets)
		merged.Customs = append(parent.Customs, child.Customs...)
		// Defaults: child overrides if set
		if child.Defaults != nil {
			merged.Defaults = child.Defaults
		}
		return &merged
	}

	// Recursively normalize all services
	var normalizeServices func(parentDefaults *symbols.ServiceDefaults, parentPolicies *symbols.PoliciesBlock, parentRateLimits *symbols.RateLimitsBlock, block *symbols.ServicesBlock)
	normalizeServices = func(parentDefaults *symbols.ServiceDefaults, parentPolicies *symbols.PoliciesBlock, parentRateLimits *symbols.RateLimitsBlock, block *symbols.ServicesBlock) {
		if block == nil {
			return
		}
		block.Defaults = MergeServiceDefaults(parentDefaults, block.Defaults)
		def.Policies = mergePolicies(parentPolicies, def.Policies)
		def.RateLimits = mergeRateLimits(parentRateLimits, def.RateLimits)
		for i := range block.Services {
			svc := &block.Services[i]
			// Merge defaults, policies, rate_limits for this service
			svcDefaults := MergeServiceDefaults(block.Defaults, nil)
			svcPolicies := mergePolicies(def.Policies, nil)
			svcRateLimits := mergeRateLimits(def.RateLimits, nil)
			// Recurse into nested services
			if len(svc.Services) > 0 {
				nestedBlock := &symbols.ServicesBlock{Defaults: svc.Defaults, Services: svc.Services}
				normalizeServices(svcDefaults, svcPolicies, svcRateLimits, nestedBlock)
				// After recursion, update nested services
				svc.Services = nestedBlock.Services
			} else {
				svc.Defaults = svcDefaults
			}
		}
	}

	normalizeServices(nil, def.Policies, def.RateLimits, def.Services)
	return def
}
