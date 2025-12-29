package ir

type PathSegmentKind string

const (
	SegmentStatic   PathSegmentKind = "static"
	SegmentParam    PathSegmentKind = "param"
	SegmentWildcard PathSegmentKind = "wildcard"
	SegmentCatchAll PathSegmentKind = "catch_all"
	SegmentRegex    PathSegmentKind = "regex"
	SegmentOptional PathSegmentKind = "optional"
)

type PathSegment struct {
	Kind    PathSegmentKind
	Name    string // param name, wildcard name
	Literal string // static literal
	Regex   string
}

type IRRoute struct {
	ID       string        `json:"id"`
	Method   string        `json:"method"`
	Path     string        `json:"path"`
	Segments []PathSegment `json:"segments,omitempty"`

	Service   string `json:"service,omitempty"`
	Operation string `json:"operation"`

	Middlewares []string `json:"middlewares,omitempty"`

	// request-time enforcement
	RequestPolicies []string `json:"request_policies,omitempty"`

	// always applied
	BaseRateLimits []string `json:"base_rate_limits,omitempty"`

	// conditional (policy-dependent)
	PolicyRateLimits []IRPolicyRateLimit `json:"policy_rate_limits,omitempty"`

	// post-resource
	ResourcePolicies []string `json:"resource_policies,omitempty"`
}

type IRPolicyRateLimit struct {
	Policy string `json:"policy"`     // policy name
	Rate   string `json:"rate_limit"` // rate limit name
}

type IRRoutes map[string]IRRoute
