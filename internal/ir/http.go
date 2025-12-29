package ir

type IRHttpConfig struct {
	BasePath         string   `json:"base_path"`
	Cors             *bool    `json:"cors,omitempty"`
	AllowedOrigins   []string `json:"allowed_origins,omitempty"`
	AllowedMethods   []string `json:"allowed_methods,omitempty"`
	AllowedHeaders   []string `json:"allowed_headers,omitempty"`
	ExposeHeaders    []string `json:"expose_headers,omitempty"`
	AllowCredentials *bool    `json:"allow_credentials,omitempty"`
	MaxAge           *int     `json:"max_age,omitempty"`
	CacheControl     string   `json:"cache_control,omitempty"`
}
