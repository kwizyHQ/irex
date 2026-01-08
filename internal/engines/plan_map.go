package engines

import (
	nodets "github.com/kwizyHQ/irex/internal/engines/node-ts"
	. "github.com/kwizyHQ/irex/internal/plan"
)

func GetPlan(ctx *PlanContext, planID string, planType string) Step {

	plansMap := map[string]map[string]Plan{
		"node-ts": {
			"build": *nodets.NodeTSWatchPlan(ctx),
			// "watch": NodeTSWatchPlan,
		},
		// Additional engines can be added here
	}

	if runtimePlans, ok := plansMap[planID]; ok {
		if selected, ok := runtimePlans[planType]; ok {
			return &PlanStep{
				Plan: &selected,
			}
		} else {
			return &PlanError{
				Message: "Plan type '" + planType + "' not found for runtime '" + planID + "'",
			}
		}
	} else {
		return &PlanError{
			Message: "Runtime '" + planID + "' not found",
		}
	}
}
