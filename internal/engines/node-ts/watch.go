package nodets

import (
	. "github.com/kwizyHQ/irex/internal/plan"
	. "github.com/kwizyHQ/irex/internal/plan/steps"
)

func NodeTSWatchPlan(*PlanContext) *Plan {
	return &Plan{
		Name: "Node TypeScript Watch",
		ID:   "watch-node-ts",
		Steps: []Step{
			&CompileTemplatesStep{
				TemplateDir: "internal/engines/node-ts/watch/templates",
			},
			&CommandStep{
				Args: []string{"echo", "hello from node-ts watch"},
			},
			// Additional steps for watch mode can be added here
		},
	}
}
