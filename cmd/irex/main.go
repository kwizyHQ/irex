package main

import (
	"log/slog"
	"os"

	"github.com/dotenv-org/godotenvvault"
	formatCmd "github.com/kwizyHQ/irex/internal/cli/common/format"
	initcmd "github.com/kwizyHQ/irex/internal/cli/common/init"
	validateCmd "github.com/kwizyHQ/irex/internal/cli/common/validate"
	"github.com/kwizyHQ/irex/internal/cli/common/watch"
	"github.com/kwizyHQ/irex/lsp"
	"github.com/spf13/cobra"
)

func main() {
	err := godotenvvault.Load()
	if err != nil {
		// println("No env file found")
	}

	// setup logging (moved to logging.go)
	SetupLogging()

	var rootCmd = &cobra.Command{
		Use:   "irex",
		Short: "IREX development CLI",
		Long:  `IREX development CLI for initializing and managing dev projects.`,
	}

	initCmd := initcmd.Run()
	watchCmd := watch.Run()

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(formatCmd.Run())
	rootCmd.AddCommand(validateCmd.NewValidateCmd())
	rootCmd.AddCommand(lsp.Run())

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
