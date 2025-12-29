package ir

type MiddlewareStage string

const (
	MiddlewarePre   MiddlewareStage = "pre"
	MiddlewarePost  MiddlewareStage = "post"
	MiddlewareError MiddlewareStage = "error"
)

type IRMiddleware struct {
	Name    string                 `json:"name"`
	Stage   MiddlewareStage        `json:"stage"`
	Handler string                 `json:"handler"`
	Order   int                    `json:"order,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type IRMiddlewares map[string]IRMiddleware
