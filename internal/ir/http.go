package ir

type IRHttpConfig struct {
	BasePath         string   `json:"base_path"`
	Cors             bool     `json:"cors"`
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	ExposeHeaders    []string `json:"expose_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
	CacheControl     string   `json:"cache_control"`
}
