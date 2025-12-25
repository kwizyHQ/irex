package pipeline

import (
	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/ir"
)

type BuildContext struct {
	ConfigAST   *ast.ConfigAST
	SchemaAST   *ast.SchemaAST
	ServicesAST *ast.ServicesAST
	// add more ASTs as needed
	ir *ir.IRBundle
}

type BuildOptions struct {
	ConfigPath string
}
