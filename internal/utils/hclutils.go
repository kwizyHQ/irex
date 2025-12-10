package utils

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// HCL attribute extraction helpers for reuse
func GetString(attr *hclsyntax.Attribute) string {
	return GetStringExpr(attr.Expr)
}
func GetStringExpr(expr hclsyntax.Expression) string {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		return lit.Val.AsString()
	}
	return ""
}
func GetBool(attr *hclsyntax.Attribute) bool {
	return GetBoolExpr(attr.Expr)
}
func GetBoolExpr(expr hclsyntax.Expression) bool {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		return lit.Val.True()
	}
	return false
}
func GetInt(attr *hclsyntax.Attribute) int {
	if lit, ok := attr.Expr.(*hclsyntax.LiteralValueExpr); ok {
		i64, _ := lit.Val.AsBigFloat().Int64()
		return int(i64)
	}
	return 0
}
func GetLiteral(attr *hclsyntax.Attribute) interface{} {
	if lit, ok := attr.Expr.(*hclsyntax.LiteralValueExpr); ok {
		val := lit.Val
		switch val.Type() {
		case val.Type():
			if val.Type().IsPrimitiveType() {
				switch {
				case val.Type().Equals(val.Type()):
					// Return the Go native value based on type
					if val.IsNull() {
						return nil
					}
					// For strings
					if val.Type().FriendlyName() == "string" {
						return val.AsString()
					}
					// For numbers
					if val.Type().FriendlyName() == "number" {
						f, _ := val.AsBigFloat().Float64()
						return f
					}
					// For bools
					if val.Type().FriendlyName() == "bool" {
						return val.True()
					}
				}
			}
		}
		return val
	}
	return nil
}
func GetStringSlice(expr hclsyntax.Expression) []string {
	var out []string
	if tuple, ok := expr.(*hclsyntax.TupleConsExpr); ok {
		for _, e := range tuple.Exprs {
			if lit, ok := e.(*hclsyntax.LiteralValueExpr); ok {
				out = append(out, lit.Val.AsString())
			}
		}
	}
	return out
}
