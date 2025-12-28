package ir

type IRRateLimit struct {
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	Limit      string               `json:"limit"`
	CountKey   []string             `json:"count_key,omitempty"`
	BucketSize *int                 `json:"bucket_size,omitempty"`
	RefillRate string               `json:"refill_rate,omitempty"`
	Burst      *int                 `json:"burst,omitempty"`
	Action     string               `json:"action"`
	Response   *IRRateLimitResponse `json:"response,omitempty"`
	Custom     bool                 `json:"custom,omitempty"`
}

type IRRateLimitResponse struct {
	StatusCode int               `json:"status_code"`
	Body       map[string]string `json:"body,omitempty"`
}

type IRRateLimits map[string]IRRateLimit
