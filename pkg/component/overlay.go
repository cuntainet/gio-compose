package component

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	"github.com/octohelm/gio-compose/pkg/iter"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/util/contextutil"
)

type OverlayState interface {
	OverlayStack
	Visible() bool
	Topmost() bool
}

type OverlayStack interface {
	Add(o OverlayState)
	Remove(o OverlayState)
}

var OverlayContext = contextutil.New[OverlayState](contextutil.Defaulter(func() OverlayState {
	return &overlayState{}
}))

type overlayState struct {
	visible  func() bool
	children []OverlayState
}

func (o *overlayState) Add(child OverlayState) {
	o.children = append(o.children, child)
}

func (o *overlayState) Remove(child OverlayState) {
	o.children = iter.Filter(o.children, func(c OverlayState) bool {
		return c != child
	})
}

func (o *overlayState) Visible() bool {
	if o.visible != nil {
		return o.visible()
	}
	return false
}

func (o *overlayState) Topmost() bool {
	for _, c := range o.children {
		if c.Visible() {
			return false
		}
	}
	return true
}

type Overlay struct {
	Visible bool
	Refs    []*Ref[Widget]
}

func (o Overlay) Build(c BuildContext) VNode {
	ref := UseRef(c, o.Visible)
	ref.Current = o.Visible

	parentState := OverlayContext.Extract(c)

	state := UseMemo(c, func() *overlayState {
		return &overlayState{
			visible: func() bool {
				return ref.Current
			},
		}
	}, []any{})

	UseEffect(c, func() func() {
		parentState.Add(state)
		return func() {
			parentState.Remove(state)
		}
	}, []any{})

	UseEffect(c, func() func() {
		w := renderer.WindowContext.Extract(c)

		spew.Dump(w)

		return func() {

		}
	}, []any{})

	if !o.Visible {
		return nil
	}
	return RootPortal().Children(
		Provider(func(ctx context.Context) context.Context {
			return OverlayContext.Inject(ctx, state)
		}).Children(
			c.ChildVNodes()...,
		),
	)
}
