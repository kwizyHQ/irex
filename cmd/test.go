package main

import (
	"fmt"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/ir/service"
	"github.com/kwizyHQ/irex/internal/utils"
)

func main() {
	fmt.Println("start parsing hcl file")
	path := filepath.Join("internal", "core", "ir", "service", "templates", "services.hcl")
	resp, err := service.GetJson(path)
	if err != nil {
		fmt.Println("error parsing hcl file:", err)
		return
	}
	fmt.Println("completed parsing hcl file")
	// write to output.json file in the same directory
	outputPath := filepath.Join("internal", "core", "ir", "service", "templates", "output.json")
	err = utils.WriteToFile(outputPath, []byte(resp))
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}
}
