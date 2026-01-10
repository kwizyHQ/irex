package steps

import (
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/plan"
)

type CreateFoldersStep struct {
	Folders []string
}

func (f *CreateFoldersStep) ID() string {
	return "create:folders"
}

func (f *CreateFoldersStep) Name() string {
	return "Folder Creation"
}

func (f *CreateFoldersStep) Description() string {
	return "Creates necessary folders in the target directory."
}

func (f *CreateFoldersStep) Run(ctx *plan.PlanContext) error {
	for _, folder := range f.Folders {
		path := filepath.Join(ctx.TargetDir, folder)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}
