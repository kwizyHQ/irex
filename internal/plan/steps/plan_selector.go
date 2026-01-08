package plan

import (
	"log/slog"

	. "github.com/kwizyHQ/irex/internal/plan"
)

type plansMap = map[string]func(ctx *PlanContext) *Plan

type PlanSelectorStep struct {
	PlansMap        plansMap
	Key             string
	DeferLoadingKey func(ctx *PlanContext) string // get function to get delayed value for specified key in context
}

func (s *PlanSelectorStep) ID() string {
	return "select:plan"
}

func (s *PlanSelectorStep) Name() string {
	return "Select Plan"
}

func (s *PlanSelectorStep) Description() string {
	return "Selects the appropriate plan based on the key of map."
}

func (s *PlanSelectorStep) Run(ctx *PlanContext) error {
	// Example selection logic; in practice, this would be more complex
	var planKey string
	if s.Key != "" {
		planKey = s.Key
	} else if s.DeferLoadingKey != nil {
		planKey = s.DeferLoadingKey(ctx)
	}
	if planFunc, exists := s.PlansMap[planKey]; exists {
		selectedPlan := planFunc(ctx)
		return selectedPlan.Execute(ctx)
	} else {
		slog.Error("No plan found for key: " + planKey)
	}
	return nil
}
