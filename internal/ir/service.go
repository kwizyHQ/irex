package ir

type IRService struct {
	Name  string `json:"name"`
	Model string `json:"model,omitempty"`
}

type IRServices map[string]IRService
