package main

import (
	"fmt"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/ast"
	"github.com/kwizyHQ/irex/internal/utils"
)

func main() {
	basePath := filepath.Join("temp")
	// get flag --schema, --service, --config if present and parse the respective hcl file
	fmt.Println("start parsing hcl file")
	path := filepath.Join(basePath, "schema", "models.hcl")
	parsed, err := ast.ParseToJsonCommon(path, "schema")
	if err != nil {
		fmt.Println("error parsing hcl file:", err)
		return
	}
	err = utils.WriteToFile(filepath.Join(basePath, "output.json"), []byte(parsed))
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}
}
