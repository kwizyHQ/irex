package assemble

import (
	"time"

	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/ir"
)

// resolveEnvRef extracts the name from functions.EnvRef. For now we just
// return the Name field as the resolved value. More advanced resolution
// (env lookup, secret handling) can be added later.
// func resolveEnvRef(ref *symbols.RuntimeSchemaOptions) (string, string) {
// 	if ref == nil {
// 		return "", ""
// 	}
// 	var uri, db string
// 	if ref.Operations != nil {
// 		if ref.Options.URI.Name != "" {
// 			uri = ref.Options.URI.Name
// 		}
// 		if ref.Options.DB.Name != "" {
// 			db = ref.Options.DB.Name
// 		}
// 	}
// 	return uri, db
// }

func prepareConfigIR(ctx *shared.BuildContext) error {
	if ctx == nil || ctx.ConfigAST == nil || ctx.ConfigAST.Project == nil {
		return nil
	}

	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}

	p := ctx.ConfigAST.Project

	cfg := ir.IRConfig{
		Project: ir.IRProject{
			Name:        p.Name,
			Description: p.Description,
			Version:     p.Version,
			Author:      p.Author,
			License:     p.License,
			Timezone:    p.Timezone,
		},
	}

	if p.Paths != nil {
		cfg.Paths = ir.IRPaths{
			Specifications: p.Paths.Specifications,
			Templates:      p.Paths.Templates,
			Output:         p.Paths.Output,
		}
	}

	if p.Generator != nil {
		cfg.Generator = ir.IRGenerator{
			GenerateSchema:  p.Generator.Schema,
			GenerateService: p.Generator.Service,
			DryRun:          p.Generator.DryRun,
			CleanBefore:     p.Generator.CleanBefore,
		}
	}

	if p.Runtime != nil {
		rt := ir.IRRuntime{
			Name:     p.Runtime.Name,
			Version:  p.Runtime.Version,
			Scaffold: p.Runtime.Scaffold,
		}

		if p.Runtime.Options != nil {
			rt.Options = ir.IRRuntimeOptions{
				PackageManager: p.Runtime.Options.PackageManager,
				Entry:          p.Runtime.Options.Entry,
				DevNodemon:     p.Runtime.Options.DevNodemon,
			}
		}

		if p.Runtime.Schema != nil {
			rt.Schema = ir.IRRuntimeSchema{
				Framework: p.Runtime.Schema.Framework,
				Version:   p.Runtime.Schema.Version,
			}
			if p.Runtime.Schema.Options != nil {
				// resolve env refs to plain strings (names)
				if p.Runtime.Schema.Options.URI.Name != "" {
					rt.Schema.Database.URI = p.Runtime.Schema.Options.URI.Name
				}
				if p.Runtime.Schema.Options.DB.Name != "" {
					rt.Schema.Database.DB = p.Runtime.Schema.Options.DB.Name
				}
			}
		}

		if p.Runtime.Service != nil {
			rt.Service = ir.IRRuntimeService{
				Framework: p.Runtime.Service.Framework,
				Version:   p.Runtime.Service.Version,
			}
			if p.Runtime.Service.Options != nil {
				rt.Service.Server = ir.IRServerConfig{
					Logger: p.Runtime.Service.Options.Logger,
					Port:   p.Runtime.Service.Options.Port,
					Host:   p.Runtime.Service.Options.Host,
				}
			}
		}

		cfg.Runtime = rt
	}

	if p.Meta != nil {
		cfg.Meta = ir.IRMeta{
			CreatedAt:        parseTimeOrEmpty(p.Meta.CreatedAt),
			GeneratorVersion: p.Meta.GeneratorVersion,
		}
	}

	ctx.IR.Config = cfg
	return nil
}

// parseTimeOrEmpty parses an RFC3339 time string or returns zero time.
func parseTimeOrEmpty(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
