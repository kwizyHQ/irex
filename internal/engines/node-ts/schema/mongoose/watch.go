package mongoose

import (
	"embed"
	"io/fs"

	"github.com/kwizyHQ/irex/internal/plan"
	steps "github.com/kwizyHQ/irex/internal/plan/steps"
)

//go:embed *
var templatesFS embed.FS

type IndexDataProvider struct{}

func (p *IndexDataProvider) DataKey() string {
	return "schema:index"
}

func (p *IndexDataProvider) Resolve(ctx *plan.PlanContext) (any, steps.Cardinality) {
	return ctx.IR, steps.Single
}

type ModelDataProvider struct{}

func (p *ModelDataProvider) DataKey() string {
	return "schema:model"
}
func (p *ModelDataProvider) Resolve(ctx *plan.PlanContext) (any, steps.Cardinality) {
	// get model values from Models
	models := make([]any, 0)
	for _, model := range ctx.IR.Models {
		models = append(models, model)
	}
	return models, steps.Many
}

func MongooseTSWatchPlan(ctx *plan.PlanContext) *plan.Plan {
	fsub, _ := fs.Sub(templatesFS, "templates")
	return &plan.Plan{
		Name: "Mongoose TypeScript Watch",
		ID:   "watch:mongoose-ts",
		Steps: []plan.Step{
			&steps.CompileTemplatesStep{
				Fs:            fsub,
				FrameworkType: plan.TemplateTypeSchema,
				FrameworkName: "mongoose",
			},
			&steps.RenderTemplatesStep{
				TemplateType: plan.TemplateTypeSchema,
				Providers: []steps.DataProvider{
					&IndexDataProvider{},
					&ModelDataProvider{},
				},
			},
		},
	}
}
