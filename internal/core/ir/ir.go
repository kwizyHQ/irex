package ir

// Package ir contains the intermediate representation (IR) structs
// used by generator engines. The IR is a fully-resolved, validated
// representation of the services/policies/rate-limits/etc. described
// by the user's DSL (HCL) and merged with defaults.

// Config is the top-level object produced after parsing, validation
// and defaults-merge. Generators should consume this structure.
type Config struct {
	// Global service settings
	BasePath         string   `json:"base_path"`
	Cors             bool     `json:"cors"`
	AllowedOrigins   []string `json:"allowed_origins,omitempty"`
	AllowedMethods   []string `json:"allowed_methods,omitempty"`
	AllowedHeaders   []string `json:"allowed_headers,omitempty"`
	ExposeHeaders    []string `json:"expose_headers,omitempty"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           *int     `json:"max_age,omitempty"`
	CacheControl     string   `json:"cache_control,omitempty"`

	// Resolved defaults applied to services and operations
	Defaults *ServiceDefaults `json:"defaults,omitempty"`

	// Policies keyed by name (includes presets and customs)
	RequestPolicies  map[string]RequestPolicy  `json:"request_policies,omitempty"`
	ResourcePolicies map[string]ResourcePolicy `json:"resource_policies,omitempty"`

	// Rate limit presets and custom definitions
	RateLimits map[string]RateLimit `json:"rate_limits,omitempty"`

	// Middlewares by name (lightweight representation)
	Middlewares map[string]Middleware `json:"middlewares,omitempty"`

	// Services tree (top-level services)
	Services []Service `json:"services,omitempty"`

	// Operations defined at global level (non-model endpoints)
	Operations []Operation `json:"operations,omitempty"`
}
