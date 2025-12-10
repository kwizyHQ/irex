package schema

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func parseFieldsBlock(body *hclsyntax.Body) ([]Field, error) {
	out := make([]Field, 0)
	for _, b := range body.Blocks {
		// each b is a field block; name can be a label or the block type
		var name string
		if len(b.Labels) > 0 {
			name = b.Labels[0]
		} else {
			name = b.Type
		}
		f := Field{Name: name}
		// attributes
		for an, attr := range b.Body.Attributes {
			v, diags := attr.Expr.Value(&hcl.EvalContext{})
			if diags.HasErrors() {
				return nil, fmt.Errorf("evaluate attribute %s: %s", an, diags.Error())
			}
			switch an {
			case "type":
				if v.Type() == cty.String {
					f.Type = v.AsString()
				}
			case "required":
				if v.True() {
					f.Required = true
				}
			case "optional":
				if v.True() {
					f.Optional = true
				}
			case "unique":
				if v.True() {
					f.Unique = true
				}
			case "default":
				var gv interface{}
				if err := gocty.FromCtyValue(v, &gv); err == nil {
					f.Default = gv
				}
			case "match":
				if v.Type() == cty.String {
					f.Match = v.AsString()
				}
			case "trim":
				if v.True() {
					f.Trim = true
				}
			case "visibility":
				if v.Type() == cty.String {
					f.Visibility = v.AsString()
				}
			case "description":
				if v.Type() == cty.String {
					f.Description = v.AsString()
				}
			}
		}

		// nested fields
		for _, nb := range b.Body.Blocks {
			if nb.Type == "fields" {
				nested, err := parseFieldsBlock(nb.Body)
				if err != nil {
					return nil, err
				}
				f.Fields = nested
			}
		}

		out = append(out, f)
	}
	return out, nil
}
