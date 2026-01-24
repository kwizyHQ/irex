package lsp

import (
	"context"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kwizyHQ/irex/internal/core/pipeline"
	"github.com/sourcegraph/jsonrpc2"
)

// simple heuristic: report a diagnostic for any line containing "TODO" or "error"
func computeDiagnostics(text string, uri string) []Diagnostic {
	var out = make([]Diagnostic, 0)
	filename, _ := UriToPath(uri)
	diags := pipeline.GetDiagnosticsForFile(filename, text)
	for _, d := range diags {
		out = append(out, Diagnostic{
			Range: Range{
				Start: Position{Line: d.Range.Start.Line, Character: d.Range.Start.Column},
				End:   Position{Line: d.Range.End.Line, Character: d.Range.End.Column},
			},
			Severity: int(d.Severity),
			Source:   "irex-lsp",
			Message:  d.Message,
			Code:     d.Code,
		})
	}
	return out
}

func publishDiagnostics(ctx context.Context, conn *jsonrpc2.Conn, uri string, diags []Diagnostic) error {
	params := PublishDiagnosticsParams{URI: uri, Diagnostics: diags}
	return conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}

// UriToPath converts an LSP file URI to a native OS file path.
func UriToPath(uriStr string) (string, error) {
	u, err := url.Parse(uriStr)
	if err != nil {
		return "", err
	}

	// 1. Unescape characters (e.g., %20 -> space, %3A -> :)
	path, err := url.PathUnescape(u.Path)
	if err != nil {
		return "", err
	}

	// 2. Handle Windows-specific leading slash
	// u.Path usually returns "/C:/..." or "/D:/..."
	if runtime.GOOS == "windows" {
		if len(path) > 2 && path[0] == '/' && path[2] == ':' {
			path = path[1:]
		}
	}

	// 3. Use filepath.FromSlash to ensure "\" on Windows and "/" on Unix
	return filepath.FromSlash(path), nil
}

// PathToUri converts a native OS file path back to an LSP URI.
func PathToUri(path string) string {
	// 1. Ensure forward slashes regardless of OS
	u := filepath.ToSlash(path)

	// 2. Windows needs a leading slash for the URI to be valid (e.g., /C:/...)
	if runtime.GOOS == "windows" && !strings.HasPrefix(u, "/") {
		u = "/" + u
	}

	// 3. Create a URL object with the "file" scheme
	res := &url.URL{
		Scheme: "file",
		Path:   u,
	}

	// .String() handles the final percent-encoding
	return res.String()
}
