package compose

import (
	"github.com/octohelm/gio-compose/pkg/layout"

	"github.com/octohelm/gio-compose/pkg/compose/internal"
)

type VNode = internal.VNode
type BuildContext = internal.BuildContext
type Component = internal.Component
type ComponentWrapper = internal.ComponentWrapper
type Widget = internal.Widget
type WidgetPainter = internal.WidgetPainter

func WidgetPainterFunc(layout func(gtx layout.Context) layout.Dimensions) WidgetPainter {
	return &widgetPainterFunc{
		layout: layout,
	}
}

type widgetPainterFunc struct {
	layout func(gtx layout.Context) layout.Dimensions
}

func (w *widgetPainterFunc) Layout(gtx layout.Context) layout.Dimensions {
	return w.layout(gtx)
}
