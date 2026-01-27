package plan

import (
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WatchEntry struct {
	Cmd       *exec.Cmd
	StartedAt time.Time
}

type WatchRegistry struct {
	mu      sync.Mutex
	entries map[string]*WatchEntry
}

func NewWatchRegistry() *WatchRegistry {
	return &WatchRegistry{entries: make(map[string]*WatchEntry)}
}

// StartOrRestart starts the given command or restarts it if already running.
// id is a unique identifier for the command (e.g. "npm-dev"), dir is the working directory, args is the command and its args.
// stopAndWaitForCommand kills a command process tree and waits for it to complete
func (r *WatchRegistry) stopAndWaitForCommand(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	_ = r.killProcessTree(cmd.Process.Pid)

	// IMPORTANT: wait synchronously
	_ = cmd.Wait()
}

func (r *WatchRegistry) StartOrRestart(id string, dir string, args []string) error {
	if len(args) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.entries == nil {
		r.entries = make(map[string]*WatchEntry)
	}

	// kill old process if present (ensure tree is killed)
	if e, ok := r.entries[id]; ok {
		r.stopAndWaitForCommand(e.Cmd)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// configure OS-specific process attributes (process group etc.)
	configureProcessAttributes(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	r.entries[id] = &WatchEntry{Cmd: cmd, StartedAt: time.Now()}
	return nil
}

// Stop stops a running command by id (kills the process tree).
func (r *WatchRegistry) Stop(id string) error {
	r.mu.Lock()
	e, ok := r.entries[id]
	if ok {
		delete(r.entries, id)
	}
	r.mu.Unlock()

	if !ok {
		return nil
	}

	r.stopAndWaitForCommand(e.Cmd)
	return nil
}

// Shutdown stops all running commands.
func (r *WatchRegistry) Shutdown() {
	r.mu.Lock()
	entries := r.entries
	r.entries = make(map[string]*WatchEntry)
	r.mu.Unlock()

	for _, e := range entries {
		if e != nil && e.Cmd != nil && e.Cmd.Process != nil {
			slog.Info("killing process: ", "pid", strconv.Itoa(e.Cmd.Process.Pid), "Name", strings.Join(e.Cmd.Args, ","))
		}
		r.stopAndWaitForCommand(e.Cmd)
	}
}

// killProcessTree is implemented in OS-specific files.
func (r *WatchRegistry) killProcessTree(pid int) error { return killProcessTree(pid) }
