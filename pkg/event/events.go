package event

type Events[T comparable] struct {
	handlers map[T][]Handler[T]
}

func (events *Events[T]) Empty() bool {
	return len(events.handlers) == 0
}

func (events *Events[T]) Reset() {
	events.handlers = nil
}

func (events *Events[T]) Watch(handlers ...Handler[T]) {
	if events.handlers == nil {
		events.handlers = map[T][]Handler[T]{}
	}

	for i := range handlers {
		h := handlers[i]
		if h == nil {
			continue
		}
		events.handlers[h.Type()] = append(events.handlers[h.Type()], h)
	}
}

func (events *Events[T]) Watched(types ...T) bool {
	if handlers := events.handlers; handlers != nil {
		for _, g := range types {
			if _, ok := handlers[g]; ok {
				return true
			}
		}
	}
	return false
}

func (events *Events[T]) Trigger(t T) {
	if handlers, ok := events.handlers[t]; ok {
		for i := range handlers {
			handlers[i].Handle()
		}
	}
}
