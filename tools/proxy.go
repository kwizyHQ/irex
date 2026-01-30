package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/kwizyHQ/irex/internal/platform/process"
	"github.com/kwizyHQ/irex/internal/watcher"
)

var (
	activePID int
)

func RunProxyCommand(command string, args []string) {
	runInfo, err := getRunInfoPath()
	if err != nil {
		fmt.Println("Error getting run info (Check if watcher is running):", err)
		return
	}
	err = ExecuteProxyCommand(command, args, runInfo.LogFile)
	if err != nil {
		fmt.Println("Error executing proxy command:", err)
		return
	}
	mgr := watcher.NewManager(
		[]string{runInfoPath},
		50*time.Millisecond,
		func(ctx context.Context, events []watcher.Event) error {
			// run command with args
			fmt.Println("Run info file changed, executing command:", command, "with args:", args)
			runInfo, err := getRunInfoPath()
			if err != nil {
				if events[0].Type == watcher.EventRemove {
					fmt.Println("Run info file removed, exiting proxy.")
					if activePID != 0 {
						process.KillProcessTree(activePID)
					}
					os.Exit(0)
				}
				return err
			}
			switch runInfo.Stage {
			case StageBuilding:
				// clean shutdown to active pid
				if activePID != 0 {
					fmt.Println("Killing active process with PID:", activePID)
					process.KillProcessTree(activePID)
					activePID = 0
				}
			case StageComplete:
				fmt.Println("Build complete, executing command.")
				err := ExecuteProxyCommand(command, args, runInfo.LogFile)
				if err != nil {
					fmt.Println("Exiting from proxy command:", err)
					return err
				}
			}
			return nil
		},
		false,
	)
	mgr.Run(context.Background())
}

func ExecuteProxyCommand(command string, args []string, logFilePath string) error {
	fmt.Println("Executing Proxy Command:", command, "with args:", args)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return err
	}
	cmd := exec.Command(command, args...)
	cmd.Stdout = io.MultiWriter(os.Stdout, logFile)
	cmd.Stderr = io.MultiWriter(os.Stderr, logFile)
	cmd.Start()
	activePID = cmd.Process.Pid
	process.ConfigureProcessAttributes(cmd)
	go func() {
		defer logFile.Close()
		err := cmd.Wait()
		if err != nil {
			fmt.Println("Proxy command exited with error:", err)
		}
		activePID = 0
	}()
	return cmd.Err
}
