package plan

import (
	"text/template"

	"github.com/kwizyHQ/irex/internal/core/pipeline"
	"github.com/kwizyHQ/irex/internal/ir"
	"github.com/kwizyHQ/irex/internal/tempdir"
)

type TemplateType string

const (
	TemplateTypeService TemplateType = "service"
	TemplateTypeSchema  TemplateType = "schema"
	TemplateTypeRuntime TemplateType = "runtime"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Trace(msg string)
	Debug(msg string)
}

type TemplateDefinition = pipeline.TemplateInfo

type TemplateBundle struct {
	Templates []pipeline.TemplateInfo
	Root      *template.Template
}

type RenderedTemplate struct {
	Name       string
	OutputPath string
	Content    []byte
}

type RenderSession struct {
	Files []RenderedTemplate
}

type CompiledTemplates map[TemplateType]TemplateBundle

type PlanContext struct {
	TargetDir         string
	ProjectName       string
	IR                *ir.IRBundle
	CompiledTemplates CompiledTemplates
	RenderSession     *RenderSession
	TmpDir            *tempdir.TempDir
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
	Plan *Plan
}

func (ps *PlanStep) ID() string {
	return ps.Plan.ID
}

func (ps *PlanStep) Name() string {
	return ps.Plan.Name
}

func (ps *PlanStep) Description() string {
	return ps.Plan.Description
}

func (ps *PlanStep) Run(ctx *PlanContext) error {
	return ps.Plan.Execute(ctx)
}
