package lsp

import (
	"github.com/spf13/cobra"
)

// Run returns a cobra subcommand that runs the language server over stdio.
func Run() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp",
		Short: "Run language server over stdio",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServer(cmd.Context())
		},
	}
	return cmd
}
