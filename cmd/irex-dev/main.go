package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	formatCmd "github.com/kwizyHQ/irex/internal/cli/common/format"
	initcmd "github.com/kwizyHQ/irex/internal/cli/common/init"
	"github.com/spf13/cobra"
)

func main() {
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
