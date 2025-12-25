package ir

// Operation represents a single endpoint/operation attached to a
// service or a global operation.
type Operation struct {
	Name        string       `json:"name"`
	Method      string       `json:"method,omitempty"`
	Path        string       `json:"path,omitempty"`
	Description string       `json:"description,omitempty"`
	Action      string       `json:"action,omitempty"`
	Apply       []ApplyBlock `json:"apply,omitempty"`
	Middlewares []string     `json:"middlewares,omitempty"`
	Policies    []string     `json:"policies,omitempty"`
	RateLimits  []string     `json:"rate_limits,omitempty"`
	// Resolved route/handler name for generators
	Route *Route `json:"route,omitempty"`
}

type Operations struct {
	Operations *[]Operation `json:"operations,omitempty"`
}
