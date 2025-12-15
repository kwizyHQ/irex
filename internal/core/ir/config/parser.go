package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/kwizyHQ/irex/internal/core/ir"
)

// parseHCLFile parses the HCL file at the given path into a ConfigDefinition struct
func parseHCLFile(path string) (*ConfigDefinition, error) {
	var def ConfigDefinition
	ctx := &hcl.EvalContext{
		Functions: ir.IRFunctions,
	}
	err := hclsimple.DecodeFile(path, ctx, &def)
	if err != nil {
		return nil, err
	}
	return &def, nil
}
