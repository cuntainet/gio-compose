package gesture

import (
	"gioui.org/io/pointer"
	"gioui.org/widget"
	"github.com/octohelm/gio-compose/pkg/event"

	giolayout "gioui.org/layout"
	"github.com/octohelm/gio-compose/pkg/layout"
)

type HandlersSetter interface {
	SetGestureHandlers(handlers ...Handler)
}

type Disabler interface {
	DisableGesture()
}

type FocusedChecker interface {
	Focused() bool
}

type Detector struct {
	Disabled bool

	event.Events[Gesture]

	clickable      widget.Clickable
	focusedChecker FocusedChecker

	pressed State[bool]
	hovered State[bool]
	focused State[bool]
}

type State[T comparable] struct {
	Value   T
	Changed bool
}

func (once *State[T]) Set(v T) {
	if once.Value != v {
		once.Value = v
		once.Changed = true
	}
}

func (c *Detector) DisableGesture() {
	c.Disabled = true
}

func (c *Detector) BindFocusedChecker(fc FocusedChecker) {
	c.focusedChecker = fc
}

func (c *Detector) SetGestureHandlers(handlers ...Handler) {
	c.Watch(handlers...)
}

func (c *Detector) ShouldDetectGestures(gestures ...Gesture) bool {
	return c.Watched(gestures...)
}

func (c *Detector) LayoutChild(gtx layout.Context, child giolayout.Widget) (dims layout.Dimensions) {
	if c.Events.Empty() || c.Disabled {
		return child(gtx)
	}

	dims = c.clickable.Layout(gtx, child)

	// FIXME better way to switch CursorPointer
	if c.clickable.Hovered() {
		if c.ShouldDetectGestures(Press, Tap, DoubleTap, LongPress) {
			layout.PostLayout(gtx.Ops, pointer.CursorPointer.Add)
		}
	} else {
		layout.Layout(gtx.Ops, pointer.CursorDefault.Add)
	}

	if focusedChecker := c.focusedChecker; focusedChecker != nil {
		c.focused.Set(focusedChecker.Focused())
	} else {
		c.focused.Set(c.clickable.Focused())
	}

	c.pressed.Set(c.clickable.Pressed())
	c.hovered.Set(c.clickable.Hovered())

	if c.pressed.Changed {
		if c.pressed.Value {
			c.Trigger(Press)
		} else {
			c.Trigger(Release)

			if c.focusedChecker == nil {
				// FIXME find out why focused not trigger
				c.focused.Set(false)
			}
		}
	}

	if c.clickable.Clicked() {
		c.Trigger(Tap)
	}

	if c.hovered.Changed {
		if c.hovered.Value {
			c.Trigger(Enter)
		} else {
			c.Trigger(Leave)
		}
	}

	if c.focused.Changed {
		if c.focused.Value {
			c.Trigger(Focus)
		} else {
			c.Trigger(Blur)
		}
	}

	return
}
