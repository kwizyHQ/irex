package nodets

import (
	"github.com/kwizyHQ/irex/internal/engines/node-ts/schema/mongoose"
	. "github.com/kwizyHQ/irex/internal/plan"
	. "github.com/kwizyHQ/irex/internal/plan/steps"
)

func NodeTSWatchPlan(ctx *PlanContext) *Plan {
	return &Plan{
		Name: "Node TypeScript Watch",
		ID:   "watch-node-ts",
		Steps: []Step{
			// let's select the schema framework here
			&PlanSelectorStep{
				PlansMap: map[string]func(ctx *PlanContext) *Plan{
					"mongoose": mongoose.MongooseTSWatchPlan,
				},
				Key: ctx.IR.Config.Runtime.Schema.Framework,
			},
			// let's select the service framework here
			&PlanSelectorStep{
				PlansMap: map[string]func(ctx *PlanContext) *Plan{
					// "fastify": FastifyWatchPlan,
				},
				Key: ctx.IR.Config.Runtime.Service.Framework,
			},
		},
	}
}
