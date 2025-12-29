package ir

type PolicyEffect string

const (
	PolicyAllow PolicyEffect = "allow"
	PolicyDeny  PolicyEffect = "deny"
)

type IRPolicyBase struct {
	Name        string       `json:"name"`
	Rule        string       `json:"rule"`
	Effect      PolicyEffect `json:"effect"`
	Description string       `json:"description,omitempty"`
}

type IRRequestPolicy struct {
	IRPolicyBase
}

type IRResourcePolicy struct {
	IRPolicyBase
}

type IRRequestPolicies map[string]IRRequestPolicy
type IRResourcePolicies map[string]IRResourcePolicy
