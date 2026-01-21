package lsp

import (
	"context"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type stdioReadWriteCloser struct {
	r *os.File
	w *os.File
}

func (s stdioReadWriteCloser) Read(p []byte) (int, error)  { return s.r.Read(p) }
func (s stdioReadWriteCloser) Write(p []byte) (int, error) { return s.w.Write(p) }
func (s stdioReadWriteCloser) Close() error                { return nil }

// RunServer starts an LSP-like JSON-RPC2 server over stdio.
func RunServer(ctx context.Context) error {
	// slog.Info("starting lsp server")
	stream := jsonrpc2.NewBufferedStream(stdioReadWriteCloser{r: os.Stdin, w: os.Stdout}, jsonrpc2.VSCodeObjectCodec{})
	handler := NewHandler()
	conn := jsonrpc2.NewConn(ctx, stream, handler)

	// block until context is canceled
	<-ctx.Done()
	_ = conn.Close()
	// slog.Info("lsp server stopped")
	return nil
}
