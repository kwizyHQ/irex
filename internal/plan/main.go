package plan

import (
	"github.com/kwizyHQ/irex/internal/ir"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Trace(msg string)
	Debug(msg string)
}

type PlanContext struct {
	TargetDir   string
	ProjectName string
	IR          ir.IRBundle
	// Logger      Logger
}

type Plan struct {
	ID          string
	Name        string
	Description string
	Steps       []Step
}

func (p *Plan) Execute(ctx *PlanContext) error {
	for _, step := range p.Steps {
		if err := step.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

type Step interface {
	ID() string
	Name() string
	Description() string
	Run(ctx *PlanContext) error
}

type PlanStep struct {
	plan *Plan
}

func (ps *PlanStep) ID() string {
	return ps.plan.ID
}

func (ps *PlanStep) Name() string {
	return ps.plan.Name
}

func (ps *PlanStep) Description() string {
	return ps.plan.Description
}

func (ps *PlanStep) Run(ctx *PlanContext) error {
	return ps.plan.Execute(ctx)
}
