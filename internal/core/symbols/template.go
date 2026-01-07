package symbols

type TemplateDefinition struct {
	Templates []TemplateBlock `hcl:"template,block"`
}

type TemplateBlock struct {
	Name   string `hcl:"name,label"`
	Data   string `hcl:"data,optional"`
	Output string `hcl:"output,optional"`
	Mode   string `hcl:"mode,optional"`
}
