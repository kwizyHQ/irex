# Watcher Manager

A flexible, optionated file/folder watcher for Go projects. Runs as a goroutine, supports debounce, filters, and event types.

## Features
- Watch multiple paths
- Debounce events
- Filter by glob patterns
- Select event types (create, write, remove, etc)
- Run as goroutine, clean shutdown

## Usage Example

```go
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/your/module/internal/watcher"
)

func main() {
	paths := []string{"./"}
	handler := func(ctx context.Context, events []watcher.Event) error {
		for _, ev := range events {
			fmt.Printf("%s: %s\n", ev.Type, ev.Path)
		}
		return nil
	}
	mgr := watcher.NewManager(paths, 200*time.Millisecond, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mgr.Run(ctx); err != nil {
		panic(err)
	}
}
```

## API
- `NewManager(paths []string, delay time.Duration, handler Handler) *Manager`
- `Manager.Run(ctx context.Context) error`
- `Handler func(ctx context.Context, events []Event) error`
- `Event` struct with `Path` and `Type`

See source files for more details.
