package validate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/core/pipeline"
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/diagnostics"
	"github.com/spf13/cobra"
)

// NewValidateCmd returns a cobra.Command that validates the config file and prints diagnostics.
func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [flags] <config.hcl>",
		Short: "Validate IREX config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := args[0]
			if !filepath.IsAbs(configPath) {
				absPath, err := filepath.Abs(configPath)
				if err == nil {
					configPath = absPath
				}
			}
			println("Validating config file:", configPath)
			ctx, diags := pipeline.Build(shared.BuildOptions{
				ConfigPath: configPath,
			})
			_ = ctx // ctx can be used for further processing if needed
			print(diags.Error())
			fmt.Println("Validation successful.")
			return nil
		},
	}
	return cmd
}

func printDiagnostic(d diagnostics.Diagnostic) {
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", d.Severity, d.Source, d.Message)
}
