package watcher

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar"
)

type Manager struct {
	paths    []string
	handler  Handler
	delay    time.Duration
	coalesce bool
}

func NewManager(
	paths []string,
	delay time.Duration,
	handler Handler,
	coalesce bool,
) *Manager {
	return &Manager{
		paths:    paths,
		delay:    delay,
		handler:  handler,
		coalesce: coalesce,
	}
}

func (m *Manager) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	watcher, err := NewFSWatcher()
	if err != nil {
		return err
	}

	// Expand paths with glob patterns and recursively watch directories
	expandedPaths, err := m.expandPaths(m.paths)
	if err != nil {
		return err
	}

	for _, p := range expandedPaths {
		if err := watcher.Add(p); err != nil {
			return err
		}
	}

	debouncer := NewDebouncer(m.delay)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = watcher.Run(ctx, debouncer.In())
	}()

	go func() {
		defer wg.Done()
		debouncer.Run(ctx)
	}()

	// handler loop (no goroutine leak)
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil

		case batch := <-debouncer.Out():
			if len(batch) == 0 {
				continue
			}
			if m.coalesce {
				batch = Coalesce(batch)
			}
			_ = m.handler(ctx, batch)
		}
	}
}

// expandPaths expands glob patterns and recursively discovers all directories
// to watch. Supports doublestar (**) patterns for recursive matching.
func (m *Manager) expandPaths(patterns []string) ([]string, error) {
	pathSet := make(map[string]struct{})

	for _, pattern := range patterns {
		// Check if pattern contains glob characters
		if containsGlobChars(pattern) {
			// Expand glob pattern using doublestar
			matches, err := doublestar.Glob(pattern)
			if err != nil {
				return nil, err
			}
			for _, match := range matches {
				if err := m.addPathRecursively(match, pathSet); err != nil {
					return nil, err
				}
			}
		} else {
			// Regular path - add it and its subdirectories if it's a directory
			if err := m.addPathRecursively(pattern, pathSet); err != nil {
				return nil, err
			}
		}
	}

	// Convert set to slice
	result := make([]string, 0, len(pathSet))
	for p := range pathSet {
		result = append(result, p)
	}

	return result, nil
}

// addPathRecursively adds a path and all its subdirectories to the pathSet
func (m *Manager) addPathRecursively(path string, pathSet map[string]struct{}) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		// Add the directory itself
		pathSet[path] = struct{}{}

		// Walk all subdirectories
		return filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				pathSet[subPath] = struct{}{}
			}
			return nil
		})
	} else {
		// For files, watch the parent directory
		pathSet[path] = struct{}{}
	}

	return nil
}

// containsGlobChars checks if a path contains glob pattern characters
func containsGlobChars(path string) bool {
	for _, char := range path {
		if char == '*' || char == '?' || char == '[' {
			return true
		}
	}
	return false
}
