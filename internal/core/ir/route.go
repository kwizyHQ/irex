package ir

// Route represents a resolved HTTP route that maps to a service operation
// or a global operation. It is a convenience struct used by generators
// to build routers.
type Route struct {
	Name        string   `json:"name,omitempty"`
	ServiceName string   `json:"service_name,omitempty"`
	Operation   string   `json:"operation,omitempty"`
	Method      string   `json:"method,omitempty"`
	Path        string   `json:"path,omitempty"`
	Description string   `json:"description,omitempty"`
	Action      string   `json:"action,omitempty"`
	Middlewares []string `json:"middlewares,omitempty"`
	Policies    []string `json:"policies,omitempty"`
	RateLimits  []string `json:"rate_limits,omitempty"`
}
