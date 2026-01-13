package mongoose

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/zclconf/go-cty/cty"
)

func TemplateFunctionsMap() template.FuncMap {
	return template.FuncMap{
		"ctyParse": func(c cty.Value) string {
			if !c.IsKnown() || c.IsNull() {
				return "null"
			}

			// var val any
			t := c.Type()

			switch {
			case t == cty.String:
				return c.AsString()
			case t == cty.Number:
				f, _ := c.AsBigFloat().Float64()
				return fmt.Sprintf("%g", f)
			case t == cty.Bool:
				return fmt.Sprintf("%t", c.True())
			case t.IsCollectionType() || t.IsTupleType():
				var items []string
				for it := c.ElementIterator(); it.Next(); {
					_, v := it.Element()
					// Recursively call for elements
					items = append(items, fmt.Sprintf("%v", v))
				}
				return "[" + strings.Join(items, ",") + "]"
			case t.IsObjectType():
				var pairs []string
				for name, _ := range t.AttributeTypes() {
					v := c.GetAttr(name)
					pairs = append(pairs, fmt.Sprintf("\"%s\":%v", name, v))
				}
				return "{" + strings.Join(pairs, ",") + "}"
			default:
				return ""
			}
		},
	}
}
