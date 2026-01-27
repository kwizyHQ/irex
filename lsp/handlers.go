package lsp

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

type Handler struct {
	mu   sync.Mutex
	docs map[string]string
}

func NewHandler() *Handler {
	return &Handler{docs: make(map[string]string)}
}

func (h *Handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "initialize":
		// minimal capabilities response
		result := map[string]interface{}{
			"capabilities": map[string]interface{}{
				"textDocumentSync": 1,
			},
		}
		_ = conn.Reply(ctx, req.ID, result)
	case "initialized":
		// no-op
		_ = conn.Reply(ctx, req.ID, nil)
	case "textDocument/didOpen":
		var params DidOpenTextDocumentParams
		if req.Params != nil {
			_ = json.Unmarshal(*req.Params, &params)
			h.mu.Lock()
			h.docs[params.TextDocument.URI] = params.TextDocument.Text
			h.mu.Unlock()
			go h.validateAndPublish(ctx, conn, params.TextDocument.URI)
		}
	case "textDocument/didChange":
		var params DidChangeTextDocumentParams
		if req.Params != nil {
			_ = json.Unmarshal(*req.Params, &params)
			if len(params.ContentChanges) > 0 {
				h.mu.Lock()
				h.docs[params.TextDocument.URI] = params.ContentChanges[0].Text
				h.mu.Unlock()
				go h.validateAndPublish(ctx, conn, params.TextDocument.URI)
			}
		}
	case "shutdown":
		_ = conn.Reply(ctx, req.ID, nil)
	case "exit":
		os.Exit(0)
	default:
		// ignore notifications and unknown methods silently
		// slog.Debug("unknown method", "method", req.Method)
	}
}

func (h *Handler) validateAndPublish(ctx context.Context, conn *jsonrpc2.Conn, uri string) {
	h.mu.Lock()
	text := h.docs[uri]
	h.mu.Unlock()

	diags := computeDiagnostics(text, uri)
	_ = publishDiagnostics(ctx, conn, uri, diags)
}
