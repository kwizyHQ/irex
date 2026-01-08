package watcher

import (
	"context"
	"time"
)

type Debouncer struct {
	delay time.Duration
	in    chan Event
	out   chan []Event
}

func NewDebouncer(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay: delay,
		in:    make(chan Event, 64),
		out:   make(chan []Event),
	}
}

func (d *Debouncer) In() chan<- Event {
	return d.in
}

func (d *Debouncer) Out() <-chan []Event {
	return d.out
}

func (d *Debouncer) Run(ctx context.Context) {
	defer close(d.out)

	var (
		timer  *time.Timer
		buffer []Event
	)

	flush := func() {
		if len(buffer) > 0 {
			d.out <- buffer
			buffer = nil
		}
	}

	for {
		select {
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			return

		case ev := <-d.in:
			buffer = append(buffer, ev)

			if timer == nil {
				timer = time.NewTimer(d.delay)
			} else {
				timer.Reset(d.delay)
			}

		case <-func() <-chan time.Time {
			if timer != nil {
				return timer.C
			}
			return nil
		}():
			flush()
			timer = nil
		}
	}
}
