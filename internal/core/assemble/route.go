package assemble

import (
	"fmt"

	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/ir"
)

// prepareRouteIR creates an IRRoute for the given operation and registers
// policies, rate-limits and middlewares from apply blocks when provided.
func prepareRouteIR(ctx *shared.BuildContext, operationName, method, path, service string, op *symbols.Operation) error {
	if ctx == nil {
		return nil
	}
	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.Routes == nil {
		ctx.IR.Routes = make(ir.IRRoutes)
	}

	id := fmt.Sprintf("%s:%s:%s", service, method, operationName)
	route := ir.IRRoute{
		ID:        id,
		Method:    method,
		Path:      path,
		Service:   service,
		Operation: operationName,
	}

	// collect apply blocks if present
	if op != nil {
		for _, a := range op.Apply {
			switch a.Type {
			case "policy":
				// policies applied at request-time
				if len(a.ToOperations) == 0 {
					route.RequestPolicies = append(route.RequestPolicies, a.Name)
				}
				// if rate_limits listed on apply block, attach to base rate limits
				if len(a.RateLimits) > 0 {
					route.BaseRateLimits = append(route.BaseRateLimits, a.RateLimits...)
				}
			case "rate_limit":
				// attach rate limit directly
				route.BaseRateLimits = append(route.BaseRateLimits, a.Name)
			default:
				// unknown apply types are ignored for now
			}
		}
	}

	ctx.IR.Routes[id] = route
	return nil
}
