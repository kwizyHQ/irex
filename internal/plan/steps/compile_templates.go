package plan

import . "github.com/kwizyHQ/irex/internal/plan"

type CompileTemplatesStep struct {
	TemplateDir string
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
	// Placeholder for template compilation logic
	return nil
}
