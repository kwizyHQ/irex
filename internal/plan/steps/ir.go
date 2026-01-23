package steps

import (
	"log/slog"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/pipeline"
	"github.com/kwizyHQ/irex/internal/plan"
)

type LoadIR struct {
	IRPath string
}

func (s *LoadIR) ID() string {
	return "load:ir"
}

func (s *LoadIR) Name() string {
	return "Load Intermediate Representation"
}

func (s *LoadIR) Description() string {
	return "Loads the Intermediate Representation (IR) from the specified path to the Plan Context. default path is 'irex.hcl'."
}

func (s *LoadIR) Run(ctx *plan.PlanContext) error {
	irBundle, err := pipeline.Build(pipeline.BuildOptions{
		ConfigPath: filepath.Join(ctx.TargetDir, s.IRPath),
	})
	if err.Error() != "no diagnostics" {
		slog.Error(err.Error())
		return nil
	}
	ctx.IR = irBundle
	return nil
}
