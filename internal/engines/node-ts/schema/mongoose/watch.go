package mongoose

import (
	"embed"
	"io/fs"

	. "github.com/kwizyHQ/irex/internal/plan"
	. "github.com/kwizyHQ/irex/internal/plan/steps"
)

//go:embed *
var templatesFS embed.FS

func MongooseTSWatchPlan(ctx *PlanContext) *Plan {
	fsub, _ := fs.Sub(templatesFS, "templates")
	return &Plan{
		Name: "Mongoose TypeScript Watch",
		ID:   "watch:mongoose-ts",
		Steps: []Step{
			&CompileTemplatesStep{
				Fs:            fsub,
				FrameworkType: FrameworkTypeSchema,
				FrameworkName: "mongoose2",
			},
		},
	}
}
