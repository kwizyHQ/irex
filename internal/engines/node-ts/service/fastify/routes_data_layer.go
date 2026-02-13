package fastify

import (
	"github.com/kwizyHQ/irex/internal/ir"
	"github.com/kwizyHQ/irex/internal/plan"
	steps "github.com/kwizyHQ/irex/internal/plan/steps"
)

type RoutesDataProvider struct{}

func (r *RoutesDataProvider) DataKey() string {
	return "service:routes"
}

func (r *RoutesDataProvider) Resolve(ctx *plan.PlanContext) (any, steps.Cardinality) {
	routesData := BuildRoutesDataLayer(ctx.IR)
	return routesData, steps.Many
}

type Route struct {
	ID          string
	Name        string
	Method      string
	Path        string
	HandlerName string
}

type Routes = []Route

func BuildRoutesDataLayer(irb *ir.IRBundle) []any {
	routes := make([]any, 0)
	for _, route := range irb.Routes {
		r := Route{
			ID:          route.ID,
			Name:        route.Operation,
			Method:      route.Method,
			Path:        route.Path,
			HandlerName: route.Operation,
		}
		routes = append(routes, r)
	}
	return routes
}
