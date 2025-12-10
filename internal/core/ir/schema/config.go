package schema

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func parseConfigBlock(body *hclsyntax.Body) (Config, error) {
	cfg := Config{DB: map[string]map[string]interface{}{}, Indexes: make([]Index, 0)}

	// attributes
	for name, attr := range body.Attributes {
		v, diags := attr.Expr.Value(&hcl.EvalContext{})
		if diags.HasErrors() {
			return cfg, fmt.Errorf("evaluate config attr %s: %s", name, diags.Error())
		}
		switch name {
		case "table":
			if v.Type() == cty.String {
				cfg.Table = v.AsString()
			}
		case "timestamps":
			if v.True() {
				cfg.Timestamps = true
			}
		case "idStrategy":
			if v.Type() == cty.String {
				cfg.IDStrategy = v.AsString()
			}
		case "description":
			if v.Type() == cty.String {
				cfg.Description = v.AsString()
			}
		case "strict":
			if v.True() {
				cfg.Strict = true
			}
		case "indexes":
			// indexes can be an attribute list
			var gv interface{}
			if err := gocty.FromCtyValue(v, &gv); err == nil {
				if arr, ok := gv.([]interface{}); ok {
					for _, item := range arr {
						if m, ok := item.(map[string]interface{}); ok {
							idx := Index{Options: map[string]interface{}{}}
							if flds, ok := m["fields"].([]interface{}); ok {
								for _, it := range flds {
									if s, ok := it.(string); ok {
										idx.Fields = append(idx.Fields, s)
									}
								}
							}
							if u, ok := m["unique"].(bool); ok {
								idx.Unique = u
							}
							if n, ok := m["name"].(string); ok {
								idx.Name = n
							}
							if t, ok := m["type"].(string); ok {
								idx.Type = t
							}
							cfg.Indexes = append(cfg.Indexes, idx)
						}
					}
				}
			}
		}
	}

	// blocks: indexes, db
	for _, b := range body.Blocks {
		switch b.Type {
		case "indexes":
			// indexes may be defined as attribute list or blocks
			// attempt to read as attributes inside indexes block
			for _, ib := range b.Body.Blocks {
				// each index block is anonymous - parse attributes inside
				idx := Index{Options: map[string]interface{}{}}
				for an, attr := range ib.Body.Attributes {
					v, diags := attr.Expr.Value(&hcl.EvalContext{})
					if diags.HasErrors() {
						continue
					}
					var gv interface{}
					if err := gocty.FromCtyValue(v, &gv); err == nil {
						switch an {
						case "fields":
							if s, ok := gv.([]interface{}); ok {
								for _, it := range s {
									if str, ok := it.(string); ok {
										idx.Fields = append(idx.Fields, str)
									}
								}
							}
						case "unique":
							if b, ok := gv.(bool); ok {
								idx.Unique = b
							}
						case "name":
							if s, ok := gv.(string); ok {
								idx.Name = s
							}
						case "type":
							if s, ok := gv.(string); ok {
								idx.Type = s
							}
						default:
							idx.Options[an] = gv
						}
					}
				}
				cfg.Indexes = append(cfg.Indexes, idx)
			}
		case "db":
			// db nested blocks like mongo, mysql
			for _, dbb := range b.Body.Blocks {
				name := dbb.Type
				cfg.DB[name] = map[string]interface{}{}
				// gather attributes and nested blocks as map values
				for an, attr := range dbb.Body.Attributes {
					v, diags := attr.Expr.Value(&hcl.EvalContext{})
					if diags.HasErrors() {
						continue
					}
					var gv interface{}
					if err := gocty.FromCtyValue(v, &gv); err == nil {
						cfg.DB[name][an] = gv
					}
				}
				// nested blocks inside db (objects) -> convert each to a map
				for _, nb := range dbb.Body.Blocks {
					// block like collation { locale = "en" }
					subMap := map[string]interface{}{}
					for an, attr := range nb.Body.Attributes {
						v, diags := attr.Expr.Value(&hcl.EvalContext{})
						if diags.HasErrors() {
							continue
						}
						var gv interface{}
						if err := gocty.FromCtyValue(v, &gv); err == nil {
							subMap[an] = gv
						}
					}
					cfg.DB[name][nb.Type] = subMap
				}
			}
		}
	}

	return cfg, nil
}
