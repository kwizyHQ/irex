package ir

type IROperation struct {
	Name        string `json:"name"`
	Service     string `json:"service,omitempty"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Action      string `json:"action,omitempty"`
	Description string `json:"description,omitempty"`

	// Added fields
	Generated bool   `json:"generated,omitempty"` // CRUD / batch
	Kind      string `json:"kind,omitempty"`      // crud | batch | custom
}

type IROperations map[string]IROperation
