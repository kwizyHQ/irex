package service

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	ast "github.com/kwizyHQ/irex/internal/core/ast"
)

// parseHCLFile parses the HCL file at the given path into a ServiceDefinition struct
func parseHCLFile(path string) (*ServiceDefinition, error) {
	var def ServiceDefinition
	ctx := &hcl.EvalContext{
		Functions: ast.ASTFunctions,
	}
	err := hclsimple.DecodeFile(path, ctx, &def)
	if err != nil {
		return nil, err
	}
	return &def, nil
}
