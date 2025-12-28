package ir

// Package ir contains the intermediate representation (IR) structs
// used by generator engines. The IR is a fully-resolved, validated
// representation of the services/policies/rate-limits/etc. described
// by the user's DSL (HCL) and merged with defaults.

// Config is the top-level object produced after parsing, validation
// and defaults-merge. Generators should consume this structure.
type IRBundle struct {
	Http             IRHttpConfig       `json:"http"`
	Services         IRServices         `json:"services"`
	Operations       IROperations       `json:"operations"`
	Routes           IRRoutes           `json:"routes"`
	Middlewares      IRMiddlewares      `json:"middlewares"`
	RequestPolicies  IRRequestPolicies  `json:"request_policies"`
	ResourcePolicies IRResourcePolicies `json:"resource_policies"`
	RateLimits       IRRateLimits       `json:"rate_limits"`
}
