package watcher

import (
	"context"
	"sync"
	"time"
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

	for _, p := range m.paths {
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
