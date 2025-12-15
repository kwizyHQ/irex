package service

// ServiceDefinition is the root struct for the services.hcl file
type ServiceDefinition struct {
	Policies   *PoliciesBlock   `hcl:"policies,block"`
	RateLimits *RateLimitsBlock `hcl:"rate_limits,block"`
	Services   *ServicesBlock   `hcl:"services,block"`
}

// --- POLICIES ---
type PoliciesBlock struct {
	Mode         string         `hcl:"mode,optional"`
	Precedence   string         `hcl:"precedence,optional"`
	ShortCircuit *bool          `hcl:"short_circuit,optional"`
	Presets      []PolicyPreset `hcl:"policy,block"`
	Customs      []PolicyCustom `hcl:"custom,block"`
	Groups       []PolicyGroup  `hcl:"group,block"`
}

// PolicyPreset represents a named policy (preset)
type PolicyPreset struct {
	Name        string `hcl:"name,label"`
	Effect      string `hcl:"effect,optional"`
	Scope       string `hcl:"scope,optional"`
	Rule        string `hcl:"rule,optional"`
	Description string `hcl:"description,optional"`
}

// PolicyCustom represents a custom policy block
type PolicyCustom struct {
	Name        string `hcl:"name,label"`
	Scope       string `hcl:"scope,optional"`
	Description string `hcl:"description,optional"`
}

type PolicyGroup struct {
	Name        string   `hcl:"name,label"`
	Scope       string   `hcl:"scope,optional"`
	Description string   `hcl:"description,optional"`
	Policies    []string `hcl:"policies,optional"`
}

// --- RATE LIMITS ---
type RateLimitsBlock struct {
	Defaults *RateLimitDefaults `hcl:"defaults,block"`
	Presets  []RateLimitPreset  `hcl:"preset,block"`
	Customs  []RateLimitCustom  `hcl:"custom,block"`
}

type RateLimitDefaults struct {
	Action     string             `hcl:"action,optional"`
	Type       string             `hcl:"type,optional"`
	CountKey   []string           `hcl:"count_key,optional"`
	Limit      string             `hcl:"limit,optional"`
	BucketSize *int               `hcl:"bucket_size,optional"`
	RefillRate string             `hcl:"refill_rate,optional"`
	Burst      *int               `hcl:"burst,optional"`
	Response   *RateLimitResponse `hcl:"response,block"`
}

type RateLimitResponse struct {
	StatusCode int               `hcl:"status_code,optional"`
	Body       map[string]string `hcl:"body,optional"`
}

type RateLimitPreset struct {
	Name       string             `hcl:"name,label"`
	Limit      string             `hcl:"limit,optional"`
	Type       string             `hcl:"type,optional"`
	CountKey   any                `hcl:"count_key,optional"`
	RefillRate string             `hcl:"refill_rate,optional"`
	BucketSize *int               `hcl:"bucket_size,optional"`
	Burst      *int               `hcl:"burst,optional"`
	Response   *RateLimitResponse `hcl:"response,block"`
}

type RateLimitCustom struct {
	Name string `hcl:"name,label"`
}

// --- SERVICES ---
type ServicesBlock struct {
	BasePath         string   `hcl:"base_path,optional"`
	Cors             *bool    `hcl:"cors,optional"`
	AllowedOrigins   []string `hcl:"allowed_origins,optional"`
	AllowedMethods   []string `hcl:"allowed_methods,optional"`
	AllowedHeaders   []string `hcl:"allowed_headers,optional"`
	ExposeHeaders    []string `hcl:"expose_headers,optional"`
	AllowCredentials *bool    `hcl:"allow_credentials,optional"`
	MaxAge           *int     `hcl:"max_age,optional"`
	CacheControl     string   `hcl:"cache_control,optional"`

	Defaults   *ServiceDefaults `hcl:"defaults,block"`
	Operations []Operation      `hcl:"operation,block"`
	Services   []Service        `hcl:"service,block"`
}

type ServiceDefaults struct {
	Pagination      *bool    `hcl:"pagination,optional"`
	Expose          *bool    `hcl:"expose,optional"`
	SoftDelete      *bool    `hcl:"soft_delete,optional"`
	CrudOperations  []string `hcl:"crud_operations,optional"`
	BatchOperations []string `hcl:"batch_operations,optional"`
	Middlewares     []string `hcl:"middlewares,optional"`
	Sorting         []string `hcl:"sorting,optional"`
	Filtering       []string `hcl:"filtering,optional"`
	Search          []string `hcl:"search,optional"`
}

type Operation struct {
	Name        string       `hcl:"name,label"`
	Method      string       `hcl:"method,optional"`
	Path        string       `hcl:"path,optional"`
	Description string       `hcl:"description,optional"`
	Action      string       `hcl:"action,optional"`
	Apply       []ApplyBlock `hcl:"apply,block"`
}

type Service struct {
	Name            string            `hcl:"name,label"`
	Model           string            `hcl:"model,optional"`
	Expose          *bool             `hcl:"expose,optional"`
	Path            string            `hcl:"path,optional"`
	CrudOperations  []string          `hcl:"crud_operations,optional"`
	BatchOperations []string          `hcl:"batch_operations,optional"`
	Middlewares     []string          `hcl:"middlewares,optional"`
	Policies        []string          `hcl:"policies,optional"`
	RateLimit       *ServiceRateLimit `hcl:"rate_limit,block"`
	Apply           []ApplyBlock      `hcl:"apply,block"`
	Operations      []Operation       `hcl:"operation,block"`
	Services        []Service         `hcl:"service,block"`
}

type ServiceRateLimit struct {
	Limit          string   `hcl:"limit,optional"`
	Action         string   `hcl:"action,optional"`
	ActionDuration string   `hcl:"action_duration,optional"`
	CountKey       []string `hcl:"count_key,optional"`
}

// ApplyBlock represents an apply block for policies or rate limits
type ApplyBlock struct {
	Type         string   `hcl:"type,label"` // "policy" or "rate_limit"
	Name         string   `hcl:"name,label"`
	ToOperations []string `hcl:"to_operations,optional"`
	RateLimits   []string `hcl:"rate_limits,optional"`
}
