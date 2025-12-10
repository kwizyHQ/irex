package schema

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func parseHCLFile(path string) ([]Model, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(src, path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("hcl parse: %s", diags.Error())
	}

	// prefer working with hclsyntax body for easier access
	syntaxBody, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected body type")
	}

	var out []Model
	for _, b := range syntaxBody.Blocks {
		if b.Type != "models" {
			continue
		}
		// each nested block inside models is a model block where Type == model name
		for _, mb := range b.Body.Blocks {
			// model blocks are written like: User { ... }
			name := mb.Type
			m := Model{Name: name, Fields: make([]Field, 0), Relations: make([]Relation, 0)}
			for _, sub := range mb.Body.Blocks {
				switch sub.Type {
				case "fields":
					fields, err := parseFieldsBlock(sub.Body)
					if err != nil {
						return nil, fmt.Errorf("parse fields for %s: %w", name, err)
					}
					m.Fields = fields
				case "config":
					cfg, err := parseConfigBlock(sub.Body)
					if err != nil {
						return nil, fmt.Errorf("parse config for %s: %w", name, err)
					}
					m.Config = cfg
				case "relations":
					rels, err := parseRelationsBlock(sub.Body)
					if err != nil {
						return nil, fmt.Errorf("parse relations for %s: %w", name, err)
					}
					m.Relations = rels
				}
			}
			out = append(out, m)
		}
	}

	// For robustness, also try to decode if file is HCL2 file root with attributes
	if len(out) == 0 {
		// attempt a generic decode into a map
		if syntaxFile, ok := file.Body.(*hclsyntax.Body); ok {
			_ = syntaxFile
		}
		// Not implemented: more generic decoding for other layouts
	}

	return out, nil
}
