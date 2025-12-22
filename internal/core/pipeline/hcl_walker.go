package pipeline

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// SymbolTable holds discovered attributes and blocks in an HCL file.
type SymbolTable struct {
	Attrs  map[string]*AttrSource
	Blocks map[string]*BlockSource
}

// AttrSource represents the source information for an HCL attribute.
type AttrSource struct {
	Path      string
	File      string
	DefRange  hcl.Range
	TypeRange hcl.Range
	ExprRange hcl.Range
}

// BlockSource represents the source information for an HCL block.
type BlockSource struct {
	Path      string
	File      string
	DefRange  hcl.Range
	BodyRange hcl.Range
}

// WalkHCLSymbols parses the given HCL file and returns a SymbolTable of all attributes and blocks.
func WalkHCLSymbols(filePath string) (SymbolTable, error) {
	var symbolsMap SymbolTable
	symbolsMap.Attrs = make(map[string]*AttrSource)
	symbolsMap.Blocks = make(map[string]*BlockSource)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return symbolsMap, err
	}
	configFile, parseErr := hclsyntax.ParseConfig(content, filePath, hcl.Pos{Line: 1, Column: 1})
	if parseErr != nil {
		return symbolsMap, parseErr
	}
	walkBody(configFile.Body.(*hclsyntax.Body), "", filePath, &symbolsMap)
	return symbolsMap, nil
}

func walkBody(
	body *hclsyntax.Body,
	prefix string,
	file string,
	symbols *SymbolTable,
) {
	// 1. Attributes
	for name, attr := range body.Attributes {
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		symbols.Attrs[path] = &AttrSource{
			Path:      path,
			File:      file,
			DefRange:  attr.Range(),
			TypeRange: attr.NameRange,
			ExprRange: attr.Expr.Range(),
		}
	}

	// 2. Blocks
	for _, block := range body.Blocks {
		blockPath := block.Type
		if prefix != "" {
			blockPath = prefix + "." + block.Type
		}

		for _, label := range block.Labels {
			blockPath += "." + label
		}

		symbols.Blocks[blockPath] = &BlockSource{
			Path:      blockPath,
			File:      file,
			DefRange:  block.Range(),
			BodyRange: block.Body.Range(),
		}

		walkBody(block.Body, blockPath, file, symbols)
	}
}
