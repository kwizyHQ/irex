package lsp

import (
	"context"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

// simple heuristic: report a diagnostic for any line containing "TODO" or "error"
func computeDiagnostics(text string) []Diagnostic {
	var out = make([]Diagnostic, 0)
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		if idx := strings.Index(l, "TODO"); idx != -1 {
			out = append(out, Diagnostic{
				Range:    Range{Start: Position{Line: i, Character: idx}, End: Position{Line: i, Character: idx + 4}},
				Severity: 2,
				Source:   "irex-lsp",
				Message:  "TODO found",
			})
		}
		if idx := strings.Index(strings.ToLower(l), "error"); idx != -1 {
			out = append(out, Diagnostic{
				Range:    Range{Start: Position{Line: i, Character: idx}, End: Position{Line: i, Character: idx + 5}},
				Severity: 1,
				Source:   "irex-lsp",
				Message:  "string 'error' found",
			})
		}
	}
	return out
}

func publishDiagnostics(ctx context.Context, conn *jsonrpc2.Conn, uri string, diags []Diagnostic) error {
	params := PublishDiagnosticsParams{URI: uri, Diagnostics: diags}
	return conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
