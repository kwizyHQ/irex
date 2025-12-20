package pipeline

import "github.com/kwizyHQ/irex/internal/core/ast"

type BuildContext struct {
	ConfigAST   *ast.ConfigAST
	ModelsAST   *ast.ModelsAST
	ServicesAST *ast.ServicesAST
	// add more ASTs as needed
}

type BuildOptions struct {
	ConfigPath string
}
