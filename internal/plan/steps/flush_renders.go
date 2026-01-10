package steps

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/plan"
)

type FlushRendersStep struct {
	DestDir string
}

func (s *FlushRendersStep) ID() string {
	return "write:flush_rendered_files"
}

func (s *FlushRendersStep) Name() string {
	return "Flush Rendered Files"
}

func (s *FlushRendersStep) Description() string {
	return "Writes all rendered files to disk specified destination"
}

func (s *FlushRendersStep) Run(ctx *plan.PlanContext) error {
	// Write all rendered files to disk
	if s.DestDir == "" {
		s.DestDir = ctx.TmpDir.Path()
	}
	for _, render := range ctx.RenderSession.Files {
		slog.Info("Writing rendered file", "path", render.OutputPath)
		fullPath := filepath.Join(s.DestDir, render.OutputPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directories for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(render.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}
	return nil
}
