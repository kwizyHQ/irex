package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotenv-org/godotenvvault"
	formatCmd "github.com/kwizyHQ/irex/internal/cli/common/format"
	initcmd "github.com/kwizyHQ/irex/internal/cli/common/init"
	validateCmd "github.com/kwizyHQ/irex/internal/cli/common/validate"
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

	var watchCmd = &cobra.Command{
		Use:   "watch",
		Short: "Watch mode (placeholder)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("watch command not implemented yet")
			os.Exit(0)
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(formatCmd.NewFormatCmd())
	rootCmd.AddCommand(validateCmd.NewValidateCmd())

	// Handle Ctrl+C (SIGINT) gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		fmt.Println("\nInterrupted (Ctrl+C). Exiting...")
		os.Exit(130)
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
