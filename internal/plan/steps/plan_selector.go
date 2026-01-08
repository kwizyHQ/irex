package plan

import (
	"log/slog"

	. "github.com/kwizyHQ/irex/internal/plan"
)

type plansMap = map[string]func(ctx *PlanContext) *Plan

type PlanSelectorStep struct {
	PlansMap plansMap
	GetByKey func(ctx *PlanContext) string // get function to get delayed value for specified key in context
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
	planKey := ctx.IR.Config.Runtime.Name // This could be derived from ctx or other parameters
	if planFunc, exists := s.PlansMap[planKey]; exists {
		selectedPlan := planFunc(ctx)
		slog.Info(selectedPlan.Name)
		return selectedPlan.Execute(ctx)
	}
	return nil
}
