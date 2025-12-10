package format

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewFormatCmd returns a cobra.Command that formats HCL files using the
// internal formatter. Flags mirror the underlying formatter flags.
func NewFormatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format [flags] [paths...]",
		Short: "Format HCL files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Rebuild os.Args for the formatter: program name + flags + args
			// Collect flags from cobra and append args
			var parts []string
			// include a placeholder program name
			parts = append(parts, "hclfmt")
			// transfer flags
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if f.Changed {
					// boolean flags may be represented as --flag
					if f.Value.Type() == "bool" {
						parts = append(parts, "--"+f.Name)
					} else {
						parts = append(parts, "--"+f.Name, f.Value.String())
					}
				}
			})

			// append positional args (files)
			parts = append(parts, args...)

			// set os.Args for the formatter to parse
			os.Args = parts

			// call underlying formatter
			// RunHCLFmt will exit on error; capture that behavior by deferring
			// a recover in case of an os.Exit.
			RunHCLFmt()
			return nil
		},
	}

	// Define flags similar to the formatter
	cmd.Flags().Bool("check", false, "perform a syntax check on the given files and produce diagnostics")
	cmd.Flags().BoolP("w", "w", false, "overwrite source files instead of writing to stdout")
	cmd.Flags().Bool("require-no-change", false, "return a non-zero status if any files are changed during formatting")
	cmd.Flags().Bool("version", false, "show the version number and immediately exit")

	// Make shorthand w also available as -w
	cmd.Flags().Lookup("w").NoOptDefVal = "true"

	return cmd
}

// Helper to show available command usage quickly
func PrintFormatHelp() {
	fmt.Fprintln(os.Stderr, strings.TrimSpace(`
Format HCL files

Usage: format [flags] [paths...]
`))
}
