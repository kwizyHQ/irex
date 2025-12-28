package ir

type IRRequestPolicy struct {
	Name        string `json:"name"`
	Rule        string `json:"rule"`
	Effect      string `json:"effect"` // allow | deny
	Description string `json:"description,omitempty"`

	// Deterministic tag for rate-limit binding
	Tag string `json:"tag"`
}

type IRResourcePolicy struct {
	Name        string `json:"name"`
	Rule        string `json:"rule"`
	Effect      string `json:"effect"` // allow | deny
	Description string `json:"description,omitempty"`
}

type IRRequestPolicies map[string]IRRequestPolicy
type IRResourcePolicies map[string]IRResourcePolicy
