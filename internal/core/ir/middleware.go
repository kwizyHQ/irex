package ir

// Middleware is a lightweight representation of a middleware that may be
// referenced by name in services/operations. Generators can use this
// to emit registration code or reorder middleware chains.
type Middleware struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config,omitempty"`
}
