package ir

// ResourcePolicy represents a policy that is evaluated in the context of
// a specific resource (e.g. owner checks). Resource policies typically
// reference resource fields and may be applied at service level.
type ResourcePolicy struct {
	Name        string `json:"name"`
	Effect      string `json:"effect,omitempty"` // allow | deny
	Scope       string `json:"scope,omitempty"`  // resource
	Rule        string `json:"rule,omitempty"`
	Description string `json:"description,omitempty"`

	// Resource policies should not have rate-limits attached; keep a
	// record of where they are applied for diagnostics.
	AppliedToService string   `json:"applied_to_service,omitempty"`
	AppliedToOps     []string `json:"applied_to_operations,omitempty"`
}
