package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"

	"github.com/kwizyHQ/irex/internal/watcher"
)

func WatchAndBuild() {
	// create initial run info file
	addRunInfoFile()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	mgr := watcher.NewManager(
		[]string{"cmd", "internal", "lsp"},
		50*time.Millisecond,
		func(ctx context.Context, events []watcher.Event) error {
			fmt.Println("Change in", events[0].Path, "received", "rebuilding...")
			return BuildIrex()
		},
		true,
	)

	go func() {
		BuildIrex()
		defer wg.Done()
		err := mgr.Run(ctx)
		if err != nil {
			fmt.Println("Error running manager:", err)
		}
	}()

	go func() {
		<-sigCh
		fmt.Println("Received interrupt signal, shutting down...")
		Cleanup()
		cancel() // last option to cancel context
	}()
	wg.Wait()
}

func BuildIrex() error {
	updateBuildingStage() // set the building stage
	exec.Command("go", "build", "./cmd/irex").Run()
	fmt.Println("Rebuild complete.")
	updateCompleteStage() // set the complete stage
	return nil
}

func Cleanup() {
	// delete the runInfoPath file if exists
	if runInfoPath != "" {
		err := os.Remove(runInfoPath)
		if err != nil {
			fmt.Println("Error deleting run info file:", err)
		}
	}
}
