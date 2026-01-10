package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"
)

const (
	targetCmd  = "irex"
	outputDir  = "temp/profiles"
	pprofUrl   = "http://localhost:6060/debug/pprof"
	duration   = 5 * time.Second
	updateFile = "./irex.hcl"
)

func Profile() {
	// 1. Setup Environment
	os.MkdirAll(outputDir, 0755)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 2. Start Application
	cmd := exec.CommandContext(ctx, targetCmd, "watch")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	log.Printf("🚀 Started %s (PID: %d)", targetCmd, cmd.Process.Pid)

	// 3. Capture Profiles in a loop
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		endTime := time.Now().Add(duration)
		for time.Now().Before(endTime) {
			select {
			case <-ticker.C:
				save(pprofUrl+"/profile?seconds="+duration.String(), "cpu.prof")
				save(pprofUrl+"/heap", "heap.prof")
				save(pprofUrl+"/goroutine", "goroutine.prof")
				// Simple file "touch" to trigger watchers
				os.Chtimes(updateFile, time.Now(), time.Now())
			case <-ctx.Done():
				return
			}
		}
		stop() // Signal we are done
	}()

	<-ctx.Done()
	log.Println("✅ Done. Opening Visualizers...")

	// 4. Open Web UIs (CPU and Heap)
	visualize(":8082", "cpu.prof")
	visualize(":8080", "heap.prof")
	visualize(":8081", "goroutine.prof")
}

func save(url, name string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, _ := os.Create(filepath.Join(outputDir, name))
	defer out.Close()
	out.ReadFrom(resp.Body)
}

func visualize(port, file string) {
	path := filepath.Join(outputDir, file)
	exec.Command("go", "tool", "pprof", "-http="+port, path).Start()
	fmt.Printf("📊 View %s at http://localhost%s\n", file, port)
}
