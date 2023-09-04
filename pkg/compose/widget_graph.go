package compose

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
)

func Graph(layout func(gtx layout.Context) layout.Dimensions) Widget {
	return &graphWidget{
		Element: node.Element{
			Name: "Graph",
		},
		layout: layout,
	}
}

var _ Widget = &graphWidget{}

type graphWidget struct {
	internal.WidgetComponent
	node.Element
	layout func(gtx layout.Context) layout.Dimensions
}

func (w *graphWidget) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	return true
}

func (w *graphWidget) New(ctx context.Context) internal.Widget {
	return &graphWidget{
		Element: node.Element{
			Name: w.Name,
		},
	}
}

func (w *graphWidget) Layout(gtx layout.Context) layout.Dimensions {
	return w.layout(gtx)
}
