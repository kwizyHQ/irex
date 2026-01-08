package watcher

import "context"

type EventType int

const (
	EventUnknown EventType = iota
	EventWrite
	EventCreate
	EventRemove
	EventRename
)

type Event struct {
	Path string
	Type EventType
}

type Handler func(ctx context.Context, events []Event) error
