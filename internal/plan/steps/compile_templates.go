package steps

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/flect"
	"github.com/kwizyHQ/irex/internal/core/pipeline"
	"github.com/kwizyHQ/irex/internal/plan"
)

type CompileTemplatesStep struct {
	Fs            fs.FS
	FrameworkType plan.TemplateType
	FrameworkName string
	TemplateFuncs template.FuncMap
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
	// 1. Build the expected User Path
	uConfig := ctx.IR.Config
	userHclPath := filepath.Join(uConfig.Paths.Templates, uConfig.Runtime.Name, string(s.FrameworkType), s.FrameworkName, "templates.hcl")

	var finalHclPath string

	// 2. Check for User Override, otherwise fallback to Temp/Embedded
	if _, err := os.Stat(userHclPath); err == nil {
		finalHclPath = userHclPath
	} else {
		// Use tempdir helper to extract embedded s.Fs
		tmp := ctx.TmpDir
		srcDir := filepath.Join("templates", string(s.FrameworkType), s.FrameworkName)

		// Assuming CopyFolder handles the extraction logic
		if err := tmp.CopyFolder(s.Fs, ".", srcDir); err != nil {
			return err
		}
		finalHclPath = filepath.Join(tmp.Path(), srcDir, "templates.hcl")
	}

	slog.Debug("Compiling templates", "path", finalHclPath)

	// 3. Parse HCL and Load Templates
	res, err := pipeline.BuildTemplate(pipeline.TemplateOptions{Path: finalHclPath})
	if err != nil {
		return err
	}

	// Get the absolute directory of the HCL to resolve relative template files
	baseDir, _ := filepath.Abs(filepath.Dir(finalHclPath))

	root, err := s.parseTemplates(baseDir, res.Templates)
	ctx.CompiledTemplates[s.FrameworkType] = plan.TemplateBundle{
		Templates: res.Templates,
		Root:      root,
	}
	return err
}

// define baseTemplateFuncs
func baseTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"camel":    flect.Camelize,
		"snake":    flect.Underscore,
		"title":    flect.Titleize,
		"pascal":   flect.Pascalize,
		"kebab":    flect.Dasherize,
		"plural":   flect.Pluralize,
		"singular": flect.Singularize,
	}
}

func (s *CompileTemplatesStep) mergeFuncs() template.FuncMap {
	funcs := baseTemplateFuncs()
	for k, v := range s.TemplateFuncs {
		funcs[k] = v
	}
	return funcs
}

// ---------------- Helper Functions ----------------
func (s *CompileTemplatesStep) parseTemplates(dir string, templates []pipeline.TemplateInfo) (*template.Template, error) {
	tmpl := template.New("root").Funcs(s.mergeFuncs())

	for _, t := range templates {
		fullPath := filepath.Join(dir, t.Name)
		// ParseFiles adds the file to the existing template set
		if _, err := tmpl.ParseFiles(fullPath); err != nil {
			return nil, err
		}
		if _, err := tmpl.New("output_path:" + t.Name).Parse(t.Output); err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}
