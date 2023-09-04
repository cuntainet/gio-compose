package gesture

import "github.com/octohelm/gio-compose/pkg/event"

type Gesture int

const (
	Press Gesture = iota + 10
	Release
	Tap
	DoubleTap
	LongPress
)

const (
	Enter Gesture = iota + 20
	Leave
	Hover
)

const (
	Focus Gesture = iota + 30
	Blur
)

type Handler = event.Handler[Gesture]

func OnPress(action func()) Handler {
	return event.NewHandler(Press, action)
}

func OnRelease(action func()) Handler {
	return event.NewHandler(Release, action)
}

func OnTap(action func()) Handler {
	return event.NewHandler(Tap, action)
}

func OnDoubleTap(action func()) Handler {
	return event.NewHandler(DoubleTap, action)
}

func OnLongPress(action func()) Handler {
	return event.NewHandler(LongPress, action)
}

func OnEnter(action func()) Handler {
	return event.NewHandler(Enter, action)
}

func OnLeave(action func()) Handler {
	return event.NewHandler(Leave, action)
}

func OnHover(action func()) Handler {
	return event.NewHandler(Hover, action)
}

func OnFocus(action func()) Handler {
	return event.NewHandler(Focus, action)
}

func OnBlur(action func()) Handler {
	return event.NewHandler(Blur, action)
}
