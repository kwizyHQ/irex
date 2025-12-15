package ir

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// all functions map as "StringName": function
var IRFunctions = map[string]function.Function{
	"only":    OnlyFunc,
	"except":  ExceptFunc,
	"with":    WithFunc,
	"without": WithoutFunc,
	"env":     EnvFunc,
}

// --- Implementation of the `only` function (Similar to `with`) ---

// OnlyFunc is a function.Function wrapper for only
var OnlyFunc = function.New(&function.Spec{
	// VarParam allows an arbitrary number of string arguments
	VarParam: &function.Parameter{Type: cty.String},
	// Returns a dynamic object type (the IR structure)
	Type: function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		// Prepare the list of strings to include
		includeList := make([]cty.Value, len(args))
		copy(includeList, args)

		return cty.ObjectVal(map[string]cty.Value{
			// The items passed to the function are the ones to be included
			"include": cty.ListVal(includeList),
			// We exclude nothing
			"exclude": cty.ListValEmpty(cty.String),
			// Setting false means only the 'include' list should be used, ignoring defaults
			"mergeDefaults": cty.False,
		}), nil
	},
})

// --- Implementation of the `except` function (Similar to `without`) ---

// ExceptFunc is a function.Function wrapper for except
var ExceptFunc = function.New(&function.Spec{
	VarParam: &function.Parameter{Type: cty.String},
	Type:     function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		// Prepare the list of strings to exclude
		excludeList := make([]cty.Value, len(args))
		copy(excludeList, args)

		return cty.ObjectVal(map[string]cty.Value{
			// We include nothing explicitly
			"include": cty.ListValEmpty(cty.String),
			// The items passed to the function are the ones to be excluded
			"exclude": cty.ListVal(excludeList),
			// Setting true means start with defaults, then apply the 'exclude' list
			"mergeDefaults": cty.True,
		}), nil
	},
})

// --- Implementation of the `with` function (Alias/Duplicate logic of `only`) ---

// WithFunc is a function.Function wrapper for with
var WithFunc = function.New(&function.Spec{
	VarParam: &function.Parameter{Type: cty.String},
	Type:     function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		// Prepare the list of strings to include
		includeList := make([]cty.Value, len(args))
		copy(includeList, args)

		return cty.ObjectVal(map[string]cty.Value{
			"include": cty.ListVal(includeList),
			"exclude": cty.ListValEmpty(cty.String),
			// This typically implies a replacement of the default list
			"mergeDefaults": cty.False,
		}), nil
	},
})

// --- Implementation of the `without` function (Alias/Duplicate logic of `except`) ---

// WithoutFunc is a function.Function wrapper for without
var WithoutFunc = function.New(&function.Spec{
	VarParam: &function.Parameter{Type: cty.String},
	Type:     function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		// Prepare the list of strings to exclude
		excludeList := make([]cty.Value, len(args))
		copy(excludeList, args)

		return cty.ObjectVal(map[string]cty.Value{
			"include": cty.ListValEmpty(cty.String),
			"exclude": cty.ListVal(excludeList),
			// This typically implies starting with the default list and removing these
			"mergeDefaults": cty.True,
		}), nil
	},
})

// EnvKind represents the kind of environment reference (env, var, secret)
type EnvKind string

const (
	EnvKindEnv    EnvKind = "env"
	EnvKindVar    EnvKind = "var"
	EnvKindSecret EnvKind = "secret"
)

// EnvRef represents a reference to an environment variable, variable, or secret
type EnvRef struct {
	Name string  `hcl:"name,attr" cty:"name"`
	Kind EnvKind `hcl:"kind,attr" cty:"kind"`
}

// EnvFunc is a function.Function wrapper for env(key)
var EnvFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "name",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.Object(map[string]cty.Type{
		"name": cty.String,
		"kind": cty.String,
	})),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		name := args[0].AsString()
		return cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal(name),
			"kind": cty.StringVal(string(EnvKindEnv)),
		}), nil
	},
})
