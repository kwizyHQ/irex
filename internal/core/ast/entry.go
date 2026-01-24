package ast

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/kwizyHQ/irex/internal/core/functions"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

func ParseHCL[T any](path string, def *T) diagnostics.Diagnostics {
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
	err := ParseHCL(path, def)
	if len(err) > 0 {
		return "", err
	}
	return ToJSON(def)
}

// ToJSON converts any Go value to a pretty-printed JSON string
func ToJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ParseFromHCLContent[T any](path string, content string, def *T) diagnostics.Diagnostics {
	r := diagnostics.NewReporter()
	ctx := &hcl.EvalContext{
		Functions: functions.ASTFunctions,
	}
	err := hclsimple.Decode(path, []byte(content), ctx, def)
	r.FromHCL(err)
	return r.All()
}
