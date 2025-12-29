package ir

type ServiceKind string

const (
	ServiceModel  ServiceKind = "model"
	ServiceSystem ServiceKind = "system"
	ServiceCustom ServiceKind = "custom"
)

type IRService struct {
	Name   string      `json:"name"`
	Kind   ServiceKind `json:"kind"`
	Model  string      `json:"model,omitempty"`
	Parent string      `json:"parent,omitempty"`
	Expose *bool       `json:"expose,omitempty"`
}
type IRServices map[string]IRService
