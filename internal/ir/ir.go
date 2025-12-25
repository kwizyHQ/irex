package ir

// Package ir contains the intermediate representation (IR) structs
// used by generator engines. The IR is a fully-resolved, validated
// representation of the services/policies/rate-limits/etc. described
// by the user's DSL (HCL) and merged with defaults.

// Config is the top-level object produced after parsing, validation
// and defaults-merge. Generators should consume this structure.
type IRBundle struct {
	Services         *Services       `json:"services,omitempty"`
	Middlewares      *Middlewares    `json:"middlewares,omitempty"`
	Routes           *Routes         `json:"routes,omitempty"`
	Operations       *Operations     `json:"operations,omitempty"`
	RateLimits       *RateLimits     `json:"rate_limits,omitempty"`
	RequestPolicies  *RequestPolicy  `json:"request_policies,omitempty"`
	ResourcePolicies *ResourcePolicy `json:"response_policies,omitempty"`
}
