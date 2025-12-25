package ir

// Middleware is a lightweight representation of a middleware that may be
// referenced by name in services/operations. Generators can use this
// to emit registration code or reorder middleware chains.
type Middleware struct {
	Name    string                 `json:"name"`
	Type    *string                `json:"type,omitempty"` // optional type pre, post, error
	Handler string                 `json:"handler"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type Middlewares struct {
	Middlewares *[]Middleware `json:"middlewares,omitempty"`
}
