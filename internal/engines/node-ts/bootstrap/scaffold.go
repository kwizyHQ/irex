package bootstrap

import (
	"embed"
	"io/fs"
	"os"

	. "github.com/kwizyHQ/irex/internal/plan"
	. "github.com/kwizyHQ/irex/internal/plan/steps"
)

//go:embed all:templates
var templatesFS embed.FS

// Scaffold performs project initialization for a TypeScript Node project.
// It will:
//   - create target directory structure
//   - run `npm init` or `yarn init`
//   - install base dependencies (dotenv, axios, pino)
//   - install schema/framework deps (mongoose, fastify)
//   - install devDependencies (typescript, ts-node, @types/node, @types/dotenv, nodemon)
//   - run `npx tsc --init`
//   - create src/* folders and basic files (.env.example, README.md, src/app.ts, src/vendor/server.ts)
func Scaffold() error {
	target := os.Getenv("IREX_TARGET")
	ctx := PlanContext{
		TargetDir:   target,
		ProjectName: os.Getenv("IREX_NAME"),
	}
	return NodeTsScaffold(&ctx).Execute(&ctx)
}

func NodeTsScaffold(ctx *PlanContext) *Plan {
	subFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		return &Plan{
			Name: "Node TypeScript Scaffold",
			ID:   "scaffold-node-ts",
			Steps: []Step{
				&PlanError{
					StepID:  "scaffold-init",
					Message: "failed to load embedded templates filesystem: " + err.Error(),
				},
			},
		}
	}
	return &Plan{
		Name: "Node TypeScript Scaffold",
		ID:   "scaffold-node-ts",
		Steps: []Step{
			&LoadIR{IRPath: "irex.hcl"},
			&CommandStep{
				Args: []string{"npm", "init", "-y"},
			},
			&CommandStep{
				Args: []string{"npm", "pkg", "set", "name=" + ctx.ProjectName},
			},
			&CommandStep{
				DescriptionOverride: "Install dev dependencies",
				Args: []string{"npm", "install", "-D",
					"typescript", "ts-node", "@types/node", "@types/dotenv", "nodemon",
				},
			},
			&CommandStep{
				DescriptionOverride: "Install dependencies",
				Args: []string{"npm", "install", "--save",
					"dotenv", "axios", "pino", "fastify", "mongoose",
				},
			},
			&CreateFoldersStep{
				Folders: []string{
					"src",
					"src/hooks",
					"src/middlewares",
					"src/utils",
					"src/vendor",
					"src/workflows",
				},
			},
			&CopyFilesStep{
				FS: subFS,
				FilesCopy: map[string]string{
					"scaffold/app.ts":       "src/app.ts",
					"scaffold/server.ts":    "src/vendor/server.ts",
					"scaffold/README.md":    "README.md",
					"scaffold/.env.example": ".env.example",
					"scaffold/nodemon.json": "nodemon.json",
				},
			},
			&CommandStep{
				DescriptionOverride: "Initialize TypeScript configuration",
				Args:                []string{"npx", "tsc", "--init"},
			},
		},
	}
}
