package ir

type IRMiddleware struct {
	Name    string                 `json:"name"`
	Stage   string                 `json:"stage"` // pre | post | error
	Handler string                 `json:"handler"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type IRMiddlewares map[string]IRMiddleware
