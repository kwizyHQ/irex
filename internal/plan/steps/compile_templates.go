package plan

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/plan"
	"github.com/kwizyHQ/irex/internal/tempdir"
)

type FrameworkType string

const (
	FrameworkTypeService FrameworkType = "service"
	FrameworkTypeSchema  FrameworkType = "schema"
)

type CompileTemplatesStep struct {
	Fs            fs.FS
	FrameworkType FrameworkType
	FrameworkName string
}

func (s *CompileTemplatesStep) ID() string {
	return "compile:templates"
}

func (s *CompileTemplatesStep) Name() string {
	return "Compile Templates"
}

func (s *CompileTemplatesStep) Description() string {
	return "Compiles templates from the specified directory."
}

func (s *CompileTemplatesStep) Run(ctx *plan.PlanContext) error {
	// user override folder
	uTP := ctx.IR.Config.Paths.Templates
	uR := ctx.IR.Config.Runtime.Name
	var path string
	userTemplatePath := filepath.Join(uTP, uR, string(s.FrameworkType), s.FrameworkName, "templates.hcl")
	// check if userTemplatePath exists
	if _, err := os.Stat(userTemplatePath); os.IsNotExist(err) {
		slog.Info("user override template path does not exist, using default templates")
		// create a temp directory over which we can write the templates
		dir := tempdir.Get()
		dir.CopyFolder(s.Fs, ".", filepath.Join("templates", string(s.FrameworkType), s.FrameworkName))
	} else {
		path = userTemplatePath
		slog.Info("user override template path exits")
	}
	slog.Info("Compiling templates from path", "path", path)
	return nil
}
