package event

type Handler[T comparable] interface {
	Type() T
	Handle(args ...any)
}

func NewHandler[T comparable](typ T, handle func()) Handler[T] {
	if handle == nil {
		return nil
	}

	return &handler[T]{
		typ: typ,
		handle: func(args ...any) {
			handle()
		},
	}
}

func NewHandlerWithArgs[T comparable](typ T, handle func(args ...any)) Handler[T] {
	if handle == nil {
		return nil
	}

	return &handler[T]{
		typ:    typ,
		handle: handle,
	}
}

type handler[T comparable] struct {
	typ    T
	handle func(args ...any)
}

func (h *handler[T]) Type() T {
	return h.typ
}

func (h *handler[T]) Handle(args ...any) {
	h.handle(args...)
}
