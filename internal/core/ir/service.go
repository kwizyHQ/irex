package ir

// ServiceDefaults mirrors defaults that apply to services/operations
// after merging configuration.
type ServiceDefaults struct {
	Pagination      *bool    `json:"pagination,omitempty"`
	Expose          *bool    `json:"expose,omitempty"`
	SoftDelete      *bool    `json:"soft_delete,omitempty"`
	CrudOperations  []string `json:"crud_operations,omitempty"`
	BatchOperations []string `json:"batch_operations,omitempty"`
	Middlewares     []string `json:"middlewares,omitempty"`
	Sorting         []string `json:"sorting,omitempty"`
	Filtering       []string `json:"filtering,omitempty"`
	Search          []string `json:"search,omitempty"`
}

// ApplyBlock is the resolved representation of an `apply` block that
// connects policies or rate-limits to operations.
type ApplyBlock struct {
	Type         string   `json:"type"` // "policy" or "rate_limit"
	Name         string   `json:"name"`
	ToOperations []string `json:"to_operations,omitempty"`
	RateLimits   []string `json:"rate_limits,omitempty"`
}

// Service represents a resolved service in the project. It may contain
// nested services and operations.
type Service struct {
	Name            string       `json:"name"`
	Model           string       `json:"model,omitempty"`
	Expose          *bool        `json:"expose,omitempty"`
	Path            string       `json:"path,omitempty"`
	CrudOperations  []string     `json:"crud_operations,omitempty"`
	BatchOperations []string     `json:"batch_operations,omitempty"`
	Middlewares     []string     `json:"middlewares,omitempty"`
	Policies        []string     `json:"policies,omitempty"`
	RateLimit       *RateLimit   `json:"rate_limit,omitempty"`
	Apply           []ApplyBlock `json:"apply,omitempty"`
	Operations      []Operation  `json:"operations,omitempty"`
	Services        []Service    `json:"services,omitempty"`
	// Resolved full path (including parent base_path/service path)
	FullPath string `json:"full_path,omitempty"`
}
