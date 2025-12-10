package format

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"golang.org/x/term"
)

const versionStr = "0.0.1-dev"

var (
	check       = flag.Bool("check", false, "perform a syntax check on the given files and produce diagnostics")
	reqNoChange = flag.Bool("require-no-change", false, "return a non-zero status if any files are changed during formatting")
	overwrite   = flag.Bool("w", false, "overwrite source files instead of writing to stdout")
	showVersion = flag.Bool("version", false, "show the version number and immediately exit")
)

var parser = hclparse.NewParser()
var diagWr hcl.DiagnosticWriter // initialized in init
var checkErrs = false
var changed []string

func init() {
	color := term.IsTerminal(int(os.Stderr.Fd()))
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
	}
	diagWr = hcl.NewDiagnosticTextWriter(os.Stderr, parser.Files(), uint(w), color)
}

// RunHCLFmt exposes a formatter similar to hashicorp's hclfmt command.
// It parses flags from os.Args and formats stdin or files accordingly.
func RunHCLFmt() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println(versionStr)
		return nil
	}

	if err := processFiles(); err != nil {
		return err
	}

	if checkErrs {
		return errors.New("one or more files contained errors")
	}

	if *reqNoChange {
		if len(changed) != 0 {
			return fmt.Errorf("file(s) were changed: %s", strings.Join(changed, ", "))
		}
	}

	return nil
}

func processFiles() error {
	if flag.NArg() == 0 {
		if *overwrite {
			return errors.New("error: cannot use -w without source filenames")
		}
		return processFile("<stdin>", os.Stdin)
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			return err
		case dir.IsDir():
			return fmt.Errorf("can't format directory %s", path)
		default:
			if err := processFile(path, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func processFile(fn string, in *os.File) error {
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

	if *check {
		_, diags := parser.ParseHCL(inSrc, fn)
		err = diagWr.WriteDiagnostics(diags)
		if err != nil {
			return fmt.Errorf("failed to write diagnostics: %w", err)
		}
		if diags.HasErrors() {
			checkErrs = true
			return nil
		}
	}

	outSrc := hclwrite.Format(inSrc)

	if !bytes.Equal(inSrc, outSrc) {
		changed = append(changed, fn)
		hasLocalChanges = true
	}

	if *overwrite {
		if hasLocalChanges {
			return os.WriteFile(fn, outSrc, 0644)
		}
		return nil
	}

	_, err = os.Stdout.Write(outSrc)
	return err
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: hclfmt [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}
