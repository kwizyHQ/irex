package ir

type RateLimitType string
type RateLimitAction string

const (
	RateFixedWindow   RateLimitType = "fixed_window"
	RateSlidingWindow RateLimitType = "sliding_window"
	RateTokenBucket   RateLimitType = "token_bucket"

	RateThrottle RateLimitAction = "throttle"
	RateBlock    RateLimitAction = "block"
)

type RateLimitWindow struct {
	Requests int    `json:"requests"`
	Window   string `json:"window"`
}

type IRRateLimitResponse struct {
	StatusCode int
	Body       map[string]string
}

type IRRateLimit struct {
	Name       string               `json:"name"`
	Type       RateLimitType        `json:"type"`
	Limit      RateLimitWindow      `json:"limit"`
	CountKeys  []string             `json:"count_keys,omitempty"`
	BucketSize *int                 `json:"bucket_size,omitempty"`
	RefillRate string               `json:"refill_rate,omitempty"`
	Burst      *int                 `json:"burst,omitempty"`
	Action     RateLimitAction      `json:"action"`
	Response   *IRRateLimitResponse `json:"response,omitempty"`
	Custom     bool                 `json:"custom,omitempty"`
}

type IRRateLimits map[string]IRRateLimit
