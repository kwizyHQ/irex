package mongoose

import (
	"embed"
	"io/fs"

	"github.com/kwizyHQ/irex/internal/plan"
	steps "github.com/kwizyHQ/irex/internal/plan/steps"
)

//go:embed *
var templatesFS embed.FS

func MongooseTSWatchPlan(ctx *plan.PlanContext) *plan.Plan {
	fsub, _ := fs.Sub(templatesFS, "templates")
	return &plan.Plan{
		Name: "Mongoose TypeScript Watch",
		ID:   "watch:mongoose-ts",
		Steps: []plan.Step{
			&steps.CompileTemplatesStep{
				Fs:            fsub,
				FrameworkType: steps.FrameworkTypeSchema,
				FrameworkName: "mongoose",
			},
		},
	}
}
