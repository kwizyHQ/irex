package ast

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/utils"
)

type ConfigAST = symbols.ConfigDefinition
type SchemaAST = symbols.ModelsSpec
type ServicesAST = symbols.ServiceDefinition

func parseConfigHCL(path string) (*symbols.ConfigDefinition, error) {
	var def symbols.ConfigDefinition
	ctx := &hcl.EvalContext{
		Functions: shared.ASTFunctions,
	}
	err := hclsimple.DecodeFile(path, ctx, &def)
	if err != nil {
		return nil, err
	}
	return &def, nil
}

func parseSchemaHCL(path string) (*symbols.ModelsSpec, error) {
	var def symbols.ModelsSpec
	ctx := &hcl.EvalContext{
		Functions: shared.ASTFunctions,
	}
	err := hclsimple.DecodeFile(path, ctx, &def)
	if err != nil {
		return nil, err
	}
	return &def, nil
}

func parseServiceHCL(path string) (*symbols.ServiceDefinition, error) {
	var def symbols.ServiceDefinition
	ctx := &hcl.EvalContext{
		Functions: shared.ASTFunctions,
	}
	err := hclsimple.DecodeFile(path, ctx, &def)
	if err != nil {
		return nil, err
	}
	return &def, nil
}

func ParseHCLCommon(path string, fileType string) (parsed interface{}, err error) {
	var def interface{}
	switch fileType {
	case "config":
		parsed, err := parseConfigHCL(path)
		if err != nil {
			return nil, err
		}
		def = parsed
	case "schema":
		parsed, err := parseSchemaHCL(path)
		if err != nil {
			return nil, err
		}
		def = parsed
	case "service":
		parsed, err := parseServiceHCL(path)
		if err != nil {
			return nil, err
		}
		def = parsed
	}
	return def, nil
}

func ParseToJsonCommon(path string, fileType string) (string, error) {
	var def interface{}
	var err error
	parsed, err := ParseHCLCommon(path, fileType)
	if err != nil {
		return "", err
	}
	def, err = utils.ToJSON(parsed)
	return def.(string), err
}
