package schema

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty/gocty"
)

func parseRelationsBlock(body *hclsyntax.Body) ([]Relation, error) {
	out := make([]Relation, 0)
	for _, b := range body.Blocks {
		var name string
		if len(b.Labels) > 0 {
			name = b.Labels[0]
		} else {
			name = b.Type
		}
		rel := Relation{Name: name}
		for an, attr := range b.Body.Attributes {
			v, diags := attr.Expr.Value(&hcl.EvalContext{})
			if diags.HasErrors() {
				return nil, fmt.Errorf("evaluate relation attr %s: %s", an, diags.Error())
			}
			var gv interface{}
			_ = gocty.FromCtyValue(v, &gv)
			switch an {
			case "ref":
				if s, ok := gv.(string); ok {
					rel.Ref = s
				}
			case "type":
				if s, ok := gv.(string); ok {
					rel.Type = s
				}
			case "localField":
				if s, ok := gv.(string); ok {
					rel.LocalField = s
				}
			case "foreignField":
				if s, ok := gv.(string); ok {
					rel.ForeignField = s
				}
			case "through":
				if s, ok := gv.(string); ok {
					rel.Through = s
				}
			case "onDelete":
				if s, ok := gv.(string); ok {
					rel.OnDelete = s
				}
			case "onUpdate":
				if s, ok := gv.(string); ok {
					rel.OnUpdate = s
				}
			case "embedded":
				if bval, ok := gv.(bool); ok {
					rel.Embedded = bval
				}
			}
		}
		out = append(out, rel)
	}
	return out, nil
}
