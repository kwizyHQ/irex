package main

import (
	"flag"
	"fmt"
	"os"

	initcmd "github.com/kwizyHQ/irex/internal/cli/common/init"
)

func usage() {
	fmt.Println("Usage: irex <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  init    Initialize a new project (interactive)")
	fmt.Println("  build   Build generated project (placeholder)")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		if err := initcmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "build":
		fmt.Println("build command not implemented yet")
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}
}
