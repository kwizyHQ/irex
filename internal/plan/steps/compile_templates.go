package plan

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	. "github.com/kwizyHQ/irex/internal/plan"
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

func (s *CompileTemplatesStep) Run(ctx *PlanContext) error {
	// user override folder
	uTP := ctx.IR.Config.Paths.Templates
	uR := ctx.IR.Config.Runtime.Name
	var path string
	userTemplatePath := filepath.Join(uTP, uR, string(s.FrameworkType), s.FrameworkName, "templates.hcl")
	// check if userTemplatePath exists
	if _, err := os.Stat(userTemplatePath); os.IsNotExist(err) {
		slog.Info("user override template path does not exist, using default templates")
		// create a temp directory over which we can write the templates
		tempDir, _ := s.copyTemplateDirectoryWithContents(".")
		print(tempDir)
	} else {
		path = userTemplatePath
		slog.Info("user override template path exits")
	}
	slog.Info("Compiling templates from path", "path", path)
	return nil
}

// copyTemplateDirectoryWithContents creates a temp directory and copies all contents from s.Fs (embed.FS) into it.
func (s *CompileTemplatesStep) copyTemplateDirectoryWithContents(srcDir string) (string, error) {
	tempDir, err := os.MkdirTemp("", "irex_templates_*")
	if err != nil {
		return "", err
	}
	var walkErr error
	err = fs.WalkDir(s.Fs, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			walkErr = err
			return nil
		}
		relPath, relErr := filepath.Rel(srcDir, path)
		if relErr != nil {
			walkErr = relErr
			return nil
		}
		targetPath := filepath.Join(tempDir, relPath)
		if d.IsDir() {
			if mkErr := os.MkdirAll(targetPath, 0755); mkErr != nil {
				walkErr = mkErr
			}
			return nil
		}
		fileData, readErr := fs.ReadFile(s.Fs, path)
		if readErr != nil {
			walkErr = readErr
			return nil
		}
		if mkDirErr := os.MkdirAll(filepath.Dir(targetPath), 0755); mkDirErr != nil {
			walkErr = mkDirErr
			return nil
		}
		if writeErr := os.WriteFile(targetPath, fileData, 0644); writeErr != nil {
			walkErr = writeErr
			return nil
		}
		return nil
	})
	if err != nil || walkErr != nil {
		os.RemoveAll(tempDir)
		if err != nil {
			return "", err
		}
		return "", walkErr
	}
	if err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	return tempDir, nil
}
