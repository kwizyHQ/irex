package assemble

import (
	"strconv"

	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/ir"
)

// parseSimpleLimit extracts requests and window from a limit string like "100/1m".
func parseSimpleLimit(limit string) (int, string) {
	if limit == "" {
		return 0, ""
	}
	// naive split on '/'
	for i := 0; i < len(limit); i++ {
		if limit[i] == '/' {
			left := limit[:i]
			right := limit[i+1:]
			n, err := strconv.Atoi(left)
			if err != nil {
				return 0, limit
			}
			return n, right
		}
	}
	// no slash -> try parse all as number
	n, err := strconv.Atoi(limit)
	if err != nil {
		return 0, limit
	}
	return n, ""
}

func prepareRateLimitsIR(ctx *shared.BuildContext) error {
	if ctx == nil || ctx.ServicesAST == nil {
		return nil
	}
	s := ctx.ServicesAST

	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.RateLimits == nil {
		ctx.IR.RateLimits = make(ir.IRRateLimits)
	}

	if s.RateLimits != nil {
		// presets
		for _, p := range s.RateLimits.Presets {
			name := p.Name
			requests, window := parseSimpleLimit(p.Limit)
			rl := ir.IRRateLimit{
				Name: name,
				Type: ir.RateFixedWindow,
				Limit: ir.RateLimitWindow{
					Requests: requests,
					Window:   window,
				},
				CountKeys:  p.CountKey,
				BucketSize: p.BucketSize,
				RefillRate: p.RefillRate,
				Burst:      p.Burst,
				Action:     ir.RateThrottle,
			}
			if p.Response != nil {
				rl.Response = &ir.IRRateLimitResponse{
					StatusCode: p.Response.StatusCode,
					Body:       p.Response.Body,
				}
			}
			ctx.IR.RateLimits[name] = rl
		}

		// customs
		for _, c := range s.RateLimits.Customs {
			name := c.Name
			rl := ir.IRRateLimit{
				Name:   name,
				Custom: true,
			}
			ctx.IR.RateLimits[name] = rl
		}
	}

	return nil
}
