package pipeline

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/ir"
)

type BuildContext struct {
	ConfigAST   *symbols.ConfigDefinition
	SchemaAST   *symbols.ModelsSpec
	ServicesAST *symbols.ServiceDefinition
	// add more ASTs as needed
	ir *ir.IRBundle
}

type BuildOptions struct {
	ConfigPath string
}
