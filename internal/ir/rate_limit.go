package ir

// RateLimitResponse represents a response to be returned when rate limit
// thresholds are exceeded.
type RateLimitResponse struct {
	StatusCode int               `json:"status_code,omitempty"`
	Body       map[string]string `json:"body,omitempty"`
}

// RateLimit models token-bucket/fixed-window/sliding-window presets and
// custom rate-limit configurations.
type RateLimit struct {
	Name       string             `json:"name"`
	Type       string             `json:"type,omitempty"` // fixed_window, sliding_window, token_bucket
	Limit      string             `json:"limit,omitempty"`
	CountKey   []string           `json:"count_key,omitempty"`
	BucketSize *int               `json:"bucket_size,omitempty"`
	RefillRate string             `json:"refill_rate,omitempty"`
	Burst      *int               `json:"burst,omitempty"`
	Action     string             `json:"action,omitempty"`
	Response   *RateLimitResponse `json:"response,omitempty"`
}

type CustomRateLimit struct {
	Name    string `json:"name"`
	Handler string `json:"handler"`
}

type RateLimits struct {
	Presets *[]RateLimit       `json:"presets,omitempty"`
	Customs *[]CustomRateLimit `json:"customs,omitempty"`
}
