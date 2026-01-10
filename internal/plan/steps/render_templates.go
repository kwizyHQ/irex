package steps

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/kwizyHQ/irex/internal/plan"
)

type Cardinality int

const (
	Single Cardinality = iota + 1
	Many
)

type DataProvider interface {
	DataKey() string
	Resolve(ctx *plan.PlanContext) (any, Cardinality)
}

type RenderTemplatesStep struct {
	TemplateType plan.TemplateType
	Providers    []DataProvider
}

func (s *RenderTemplatesStep) ID() string {
	return "render:templates"
}

func (s *RenderTemplatesStep) Name() string {
	return "Render Templates"
}

func (s *RenderTemplatesStep) Description() string {
	return "Renders templates using IR-driven data providers."
}

func (s *RenderTemplatesStep) Run(ctx *plan.PlanContext) error {
	if ctx.IR == nil {
		return fmt.Errorf("IR not loaded")
	}

	bundle, ok := ctx.CompiledTemplates[s.TemplateType]
	if !ok {
		// If it's missing, we can just skip this step or log a warning
		slog.Warn("Skipping service generation: No service templates found")
		return nil
	}

	for _, provider := range s.Providers {
		data, card := provider.Resolve(ctx)
		dataKey := provider.DataKey()
		switch card {
		case Many:
			for _, cT := range bundle.Templates {
				if cT.Data == dataKey {
					slog.Debug("Executing template: " + cT.Name)
					for _, item := range data.([]any) {
						slog.Debug("Rendering item for template: " + cT.Name)
						rt, err := renderTemplate(bundle, cT, item)
						if err != nil {
							return err
						} else {
							ctx.RenderSession.Files = append(ctx.RenderSession.Files, rt)
						}
					}
				}
			}
		case Single:
			for _, cT := range bundle.Templates {
				if cT.Data == dataKey {
					slog.Debug("Executing template: " + cT.Name)
					rt, err := renderTemplate(bundle, cT, data)
					if err != nil {
						return err
					} else {
						ctx.RenderSession.Files = append(ctx.RenderSession.Files, rt)
					}
				}
			}
		}
	}

	return nil
}

func renderTemplate(bundle plan.TemplateBundle, t plan.TemplateDefinition, templateData any) (plan.RenderedTemplate, error) {
	var templateBuf bytes.Buffer
	var outputPathBuf bytes.Buffer
	err := bundle.Root.ExecuteTemplate(&templateBuf, t.Name, templateData)
	if err != nil {
		return plan.RenderedTemplate{}, err
	}
	err = bundle.Root.ExecuteTemplate(&outputPathBuf, "output_path:"+t.Name, templateData)
	if err != nil {
		return plan.RenderedTemplate{}, err
	}
	renderedTemplate := plan.RenderedTemplate{
		Name:       t.Name,
		OutputPath: outputPathBuf.String(),
		Content:    templateBuf.Bytes(),
	}
	return renderedTemplate, nil
}
