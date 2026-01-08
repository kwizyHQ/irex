package watcher

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

type FSWatcher struct {
	w *fsnotify.Watcher
}

func NewFSWatcher() (*FSWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &FSWatcher{w: w}, nil
}

func (fw *FSWatcher) Add(path string) error {
	return fw.w.Add(path)
}

func (fw *FSWatcher) Run(ctx context.Context, out chan<- Event) error {
	defer fw.w.Close()

	for {
		select {
		case <-ctx.Done():
			return nil

		case ev := <-fw.w.Events:
			out <- Event{
				Path: ev.Name,
				Type: mapEvent(ev),
			}

		case <-fw.w.Errors:
			// ignore or log upstream
		}
	}
}

func mapEvent(ev fsnotify.Event) EventType {
	switch {
	case ev.Op&fsnotify.Write != 0:
		return EventWrite
	case ev.Op&fsnotify.Create != 0:
		return EventCreate
	case ev.Op&fsnotify.Remove != 0:
		return EventRemove
	case ev.Op&fsnotify.Rename != 0:
		return EventRename
	default:
		return EventUnknown
	}
}
