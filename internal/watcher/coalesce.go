package watcher

func Coalesce(events []Event) []Event {
	seen := make(map[string]Event)

	for _, e := range events {
		prev, ok := seen[e.Path]
		if !ok || stronger(e.Type, prev.Type) {
			seen[e.Path] = e
		}
	}

	out := make([]Event, 0, len(seen))
	for _, e := range seen {
		out = append(out, e)
	}
	return out
}

func stronger(a, b EventType) bool {
	order := map[EventType]int{
		EventRemove:  4,
		EventRename:  3,
		EventWrite:   2,
		EventCreate:  1,
		EventUnknown: 0,
	}
	return order[a] > order[b]
}
