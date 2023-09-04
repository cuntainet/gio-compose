package internal

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
)

type Widget interface {
	node.Node
	WidgetCreator
	WidgetPatcher
	WidgetPainter
}

type WidgetCreator interface {
	New(ctx context.Context) Widget
}

type WidgetPatcher interface {
	Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool
}

type WidgetPainter interface {
	Layout(gtx layout.Context) layout.Dimensions
}
