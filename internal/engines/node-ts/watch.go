package nodets

import (
	"github.com/kwizyHQ/irex/internal/engines/node-ts/schema/mongoose"
	"github.com/kwizyHQ/irex/internal/plan"
	"github.com/kwizyHQ/irex/internal/plan/steps"
)

func NodeTSWatchPlan(ctx *plan.PlanContext) *plan.Plan {
	return &plan.Plan{
		Name: "Node TypeScript Watch",
		ID:   "watch-node-ts",
		Steps: []plan.Step{
			// let's select the schema framework here
			&steps.PlanSelectorStep{
				PlansMap: map[string]func(ctx *plan.PlanContext) *plan.Plan{
					"mongoose": mongoose.MongooseTSWatchPlan,
				},
				Key: ctx.IR.Config.Runtime.Schema.Framework,
			},
			// let's select the service framework here
			&steps.PlanSelectorStep{
				PlansMap: map[string]func(ctx *plan.PlanContext) *plan.Plan{
					// "fastify": FastifyWatchPlan,
				},
				Key: ctx.IR.Config.Runtime.Service.Framework,
			},
			&steps.FlushRendersStep{
				DestDir: ".",
			},
		},
	}
}
