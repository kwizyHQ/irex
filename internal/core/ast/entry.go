package ast

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/kwizyHQ/irex/internal/core/functions"
	"github.com/kwizyHQ/irex/internal/diagnostics"
	"github.com/kwizyHQ/irex/internal/utils"
)

func ParseHCL[T any](path string, def *T) error {
	r := diagnostics.NewReporter()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		r.Error("We couldn't find the file at "+path, diagnostics.Range{}, "config.not_found", "pipeline")
		return r.All()
	}
	ctx := &hcl.EvalContext{
		Functions: functions.ASTFunctions,
	}
	err := hclsimple.DecodeFile(path, ctx, def)
	r.FromHCL(err)
	return r.All()
}

func ParseToJson[T any](path string, def *T) (string, error) {
	err := ParseHCL(path, def).(diagnostics.Diagnostics)
	if err.HasErrors() {
		return "", err
	}
	return utils.ToJSON(def)
}
