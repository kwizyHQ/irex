package shared

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/ir"
)

type ConfigAST = symbols.ConfigDefinition
type SchemaAST = symbols.ModelsSpec
type ServicesAST = symbols.ServiceDefinition
type IRBundle = ir.IRBundle

type BuildContext struct {
	ConfigAST   *ConfigAST
	SchemaAST   *SchemaAST
	ServicesAST *ServicesAST
	// add more ASTs as needed
	IR *IRBundle
}

type BuildOptions struct {
	ConfigPath string
}
