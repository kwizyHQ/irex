package format

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewFormatCmd returns a cobra.Command that formats HCL files using the internal formatter.
func Run() *cobra.Command {
	var check, overwrite, requireNoChange bool

	cmd := &cobra.Command{
		Use:   "format [flags] [paths...]",
		Short: "Format HCL files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHCLFmt(args, check, overwrite, requireNoChange)
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "perform a syntax check on the given files and produce diagnostics")
	cmd.Flags().BoolVarP(&overwrite, "w", "w", false, "overwrite source files instead of writing to stdout")
	cmd.Flags().BoolVar(&requireNoChange, "require-no-change", false, "return a non-zero status if any files are changed during formatting")
	cmd.Flags().Lookup("w").NoOptDefVal = "true"

	return cmd
}

func runHCLFmt(paths []string, check, overwrite, requireNoChange bool) error {
	parser := hclparse.NewParser()
	color := term.IsTerminal(int(os.Stderr.Fd()))
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
	}
	diagWr := hcl.NewDiagnosticTextWriter(os.Stderr, parser.Files(), uint(w), color)
	var changed []string
	var checkErrs bool

	if len(paths) == 0 {
		fmt.Print("No paths provided. Enter HCL content (end with Ctrl+D):\n")
		if overwrite {
			return errors.New("error: cannot use -w without source filenames")
		}
		return processHclFile("<stdin>", os.Stdin, parser, diagWr, check, overwrite, &checkErrs, &changed)
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return fmt.Errorf("can't format directory %s", path)
		}
		if err := processHclFile(path, nil, parser, diagWr, check, overwrite, &checkErrs, &changed); err != nil {
			return err
		}
	}

	if checkErrs {
		return errors.New("one or more files contained errors")
	}
	if requireNoChange && len(changed) != 0 {
		return fmt.Errorf("file(s) were changed: %s", strings.Join(changed, ", "))
	}
	return nil
}

func processHclFile(fn string, in *os.File, parser *hclparse.Parser, diagWr hcl.DiagnosticWriter, check, overwrite bool, checkErrs *bool, changed *[]string) error {
	var err error
	hasLocalChanges := false
	if in == nil {
		in, err = os.Open(fn)
		if err != nil {
			return fmt.Errorf("failed to open %s: %s", fn, err)
		}
		defer in.Close()
	}

	inSrc, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("failed to read %s: %s", fn, err)
	}

	if check {
		_, diags := parser.ParseHCL(inSrc, fn)
		err = diagWr.WriteDiagnostics(diags)
		if err != nil {
			return fmt.Errorf("failed to write diagnostics: %w", err)
		}
		if diags.HasErrors() {
			*checkErrs = true
			return nil
		}
	}

	outSrc := hclwrite.Format(inSrc)

	if !bytes.Equal(inSrc, outSrc) {
		*changed = append(*changed, fn)
		hasLocalChanges = true
	}

	if overwrite {
		if hasLocalChanges {
			return os.WriteFile(fn, outSrc, 0644)
		}
		return nil
	}

	_, err = os.Stdout.Write(outSrc)
	return err
}
