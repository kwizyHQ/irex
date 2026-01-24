package pipeline

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/validate"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

func GetDiagnosticsForFile(filename string, content string) diagnostics.Diagnostics {
	r := diagnostics.NewReporter()
	fileType := GetFileType(filename)
	switch fileType {
	case "config":
		configAST := &shared.ConfigAST{}
		r.SetFilename(filename)
		r.ExtendWithFilename(ast.ParseFromHCLContent(filename, content, configAST))
		r.ExtendWithFilename(validate.ValidateConfig(configAST))
	case "schema":
		schemaAST := &shared.SchemaAST{}
		r.SetFilename(filename)
		r.ExtendWithFilename(ast.ParseFromHCLContent(filename, content, schemaAST))
		r.ExtendWithFilename(validate.ValidateSchema(schemaAST))
	case "service":
		serviceAST := &shared.ServicesAST{}
		r.SetFilename(filename)
		r.ExtendWithFilename(ast.ParseFromHCLContent(filename, content, serviceAST))
		r.ExtendWithFilename(validate.ValidateService(serviceAST))
	}
	// let's merge the ranges as well (if not zeroRange)
	diags := r.All()
	table, err := WalkHCLSymbols(filename)
	if err != nil {
		return r.All()
	}
	for i, d := range diags {
		if d.HclPath != "" {
			var rng hcl.Range
			if table.Attrs[d.HclPath] != nil {
				rng = table.Attrs[d.HclPath].ExprRange
			} else {
				// if we have no Attr, then we should issue error to parent block e.g. for project.name attribute we should point to project block
				// split the hclPath by dot and remove last segment
				parts := strings.Split(d.HclPath, ".")
				if len(parts) > 1 {
					parentPath := ""
					for _, p := range parts[:len(parts)-1] {
						if parentPath == "" {
							parentPath = p
						} else {
							parentPath += "." + p
						}
					}
					if table.Blocks[parentPath] != nil {
						rng = table.Blocks[parentPath].BodyRange
					}
				}
			}
			if table.Blocks[d.HclPath] != nil {
				rng = table.Blocks[d.HclPath].BodyRange
			}
			if rng.Empty() == false {
				diags[i].Range = diagnostics.Range{
					Start: diagnostics.Position(rng.Start),
					End:   diagnostics.Position(rng.End),
				}
			}
		}
		// let's convert to vscode style range (end is exclusive)
		if diags[i].Range.End.Line > 0 && diags[i].Range.End.Column > 0 {
			diags[i].Range.End.Line -= 1
			diags[i].Range.End.Column -= 1
		}
		// convert start range too
		if diags[i].Range.Start.Line > 0 && diags[i].Range.Start.Column > 0 {
			diags[i].Range.Start.Line -= 1
			diags[i].Range.Start.Column -= 1
		}
	}
	return diags
}

// GetFileType infers the type based on filename first, then the closest parent folder.
func GetFileType(path string) string {
	filename := filepath.Base(path)

	// 1. Check for specific filename overrides
	switch filename {
	case "irex.hcl":
		return "config"
	case "templates.hcl":
		return "template"
	}

	// 2. Traverse up the directory tree to find the closest parent
	// We use filepath.Dir in a loop to look at each parent level
	current := filepath.Dir(path)
	for {
		parentDir := filepath.Base(current)

		if parentDir == "schema" {
			return "schema"
		}
		if parentDir == "service" {
			return "service"
		}

		// Stop if we reach the root directory
		next := filepath.Dir(current)
		if next == current {
			break
		}
		current = next
	}

	return "unknown"
}
