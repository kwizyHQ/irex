package fastify

import (
	"embed"
	"io/fs"

	"github.com/kwizyHQ/irex/internal/plan"
	steps "github.com/kwizyHQ/irex/internal/plan/steps"
)

//go:embed *
var templatesFS embed.FS

type AppDataProvider struct{}

func (p *AppDataProvider) DataKey() string {
	return "service:app"
}

func (p *AppDataProvider) Resolve(ctx *plan.PlanContext) (any, steps.Cardinality) {
	appData := BuildAppDataLayer(ctx.IR)
	return appData, steps.Single
}

func FastifyTSWatchPlan(ctx *plan.PlanContext) *plan.Plan {
	fsub, _ := fs.Sub(templatesFS, "templates")
	return &plan.Plan{
		Name: "Fastify TypeScript Watch",
		ID:   "watch:fastify-ts",
		Steps: []plan.Step{
			&steps.CompileTemplatesStep{
				Fs:            fsub,
				FrameworkType: plan.TemplateTypeService,
				FrameworkName: "fastify",
				// TemplateFuncs: TemplateFunctionsMap(),
			},
			&steps.RenderTemplatesStep{
				TemplateType: plan.TemplateTypeService,
				Providers: []steps.DataProvider{
					&AppDataProvider{},
				},
			},
		},
	}
}
