package ir

// RequestPolicy is a policy that applies at request (or request-like)
// scopes. These are policies that evaluate using request-level context
// (e.g. authentication, IP, headers) and can be applied to operations.
type RequestPolicy struct {
	Name        string `json:"name"`
	Effect      string `json:"effect,omitempty"` // allow | deny
	Scope       string `json:"scope,omitempty"`  // request
	Rule        string `json:"rule,omitempty"`
	Description string `json:"description,omitempty"`

	// When the policy was applied via an `apply` block, the operations
	// this policy targets ("*" for all) — filled during IR preparation.
	ToOperations []string `json:"to_operations,omitempty"`

	// Associated rate limits referenced alongside the policy application
	RateLimits []string `json:"rate_limits,omitempty"`
}
