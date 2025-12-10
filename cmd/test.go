package main

import (
	"fmt"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/ir/schema"
)

func main() {
	fmt.Println("start parsing hcl file")
	path := filepath.Join("internal", "core", "ir", "schema", "templates", "models.hcl")
	resp, err := schema.GetJsonModels(path)
	if err != nil {
		fmt.Println("error parsing hcl file:", err)
		return
	}
	fmt.Println("completed parsing hcl file")
	fmt.Println(resp)
}
