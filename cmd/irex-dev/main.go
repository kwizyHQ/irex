package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/dotenv-org/godotenvvault"
	formatCmd "github.com/kwizyHQ/irex/internal/cli/common/format"
	initcmd "github.com/kwizyHQ/irex/internal/cli/common/init"
	validateCmd "github.com/kwizyHQ/irex/internal/cli/common/validate"
	"github.com/kwizyHQ/irex/internal/cli/common/watch"
	"github.com/spf13/cobra"
)

var startTime = time.Now()

type elapsedHandler struct{ slog.Handler }

func (h elapsedHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(
		slog.String("elapsed", time.Since(startTime).Round(time.Millisecond).String()),
	)
	return h.Handler.Handle(ctx, r)
}

func main() {

	err := godotenvvault.Load()
	if err != nil {
		println("No env file found")
	}
	// set logger defaults
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(&elapsedHandler{Handler: handler})
	// logger := slog.New(handler)
	slog.SetDefault(logger)

	var rootCmd = &cobra.Command{
		Use:   "irexd",
		Short: "IREX development CLI",
		Long:  `IREX development CLI for initializing and managing dev projects.`,
	}

	initCmd := initcmd.Run()
	watchCmd := watch.Run()

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(formatCmd.Run())
	rootCmd.AddCommand(validateCmd.NewValidateCmd())

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
