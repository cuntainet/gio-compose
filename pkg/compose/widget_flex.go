package compose

import (
	"context"
	"image"
	"math"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/unit"

	"gioui.org/op"

	giolayout "gioui.org/layout"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"github.com/octohelm/gio-compose/pkg/layout/direction"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/size"
)

func Column(modifiers ...modifier.Modifier[any]) VNode {
	return H(&flexWidget{
		Element: node.Element{
			Name: "Column",
		},
		Axis: direction.Vertical,
	}, modifiers...)
}

func Row(modifiers ...modifier.Modifier[any]) VNode {
	return H(&flexWidget{
		Element: node.Element{
			Name: "Row",
		},
		Axis: direction.Horizontal,
	}, modifiers...)
}

var _ Widget = &flexWidget{}

type flexWidget struct {
	internal.WidgetComponent
	node.Element
	Axis direction.Axis

	flexWidgetAttrs
	list *giolayout.List
}

func (fw *flexWidget) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	attrs := &flexWidgetAttrs{
		container: newContainer(),
	}

	modifier.Modify[any](attrs, modifiers...)

	return cmp.UpdateWhen(
		cmp.Not(fw.flexWidgetAttrs.Eq(attrs)),
		&fw.flexWidgetAttrs, attrs,
	)
}

type flexWidgetAttrs struct {
	*container

	layout.Spacer
	layout.Aligner
	layout.Arrangementer
	layout.Scrollable
}

func (attrs *flexWidgetAttrs) Eq(v *flexWidgetAttrs) cmp.Result {
	return cmp.All(
		attrs.container.Eq(v.container),
		attrs.Scrollable.Eq(&v.Scrollable),
		attrs.Spacer.Eq(&v.Spacer),
		attrs.Aligner.Eq(&v.Aligner),
		attrs.Arrangementer.Eq(&v.Arrangementer),
	)
}

func (fw *flexWidget) New(ctx context.Context) Widget {
	return &flexWidget{
		Element: node.Element{
			Name: fw.Name,
		},
		Axis: fw.Axis,
	}
}

func (fw *flexWidget) Layout(gtx layout.Context) layout.Dimensions {
	if fw.Scrollable.Enabled {
		return fw.container.Layout(gtx, fw, WidgetPainterFunc(func(gtx layout.Context) layout.Dimensions {
			children := make([]Widget, 0)

			for child := range node.IterChildElement(context.Background(), fw) {
				if w, ok := child.(Widget); ok {
					children = append(children, w)
				}
			}

			if fw.list == nil {
				// bind instance to store list state
				fw.list = &giolayout.List{
					Axis:      fw.Scrollable.Axis.LayoutAxis(),
					Alignment: fw.Alignment.LayoutAlignment(),
				}
			}

			maxViewSize := 0

			switch fw.Scrollable.Axis {
			case direction.Vertical:
				maxViewSize = gtx.Constraints.Max.Y
			case direction.Horizontal:
				maxViewSize = gtx.Constraints.Max.X
			}

			childrenOffsets := map[int]int{}

			isVisible := func(index int) bool {
				return index >= fw.list.Position.First && index <= fw.list.Position.First+fw.list.Position.Count
			}

			defer func() {
				for index := range children {
					positionChild(fw, children[index], func() (x, y unit.Dp) {
						childOffset := -maxViewSize

						if isVisible(index) {
							if offset, ok := childrenOffsets[index]; ok {
								childOffset = offset
							}
						} else {
							childOffset = -maxViewSize
						}

						x, y = unit.Dp(0), unit.Dp(0)

						switch fw.Scrollable.Axis {
						case direction.Vertical:
							x, y = 0, gtx.Metric.PxToDp(childOffset)
						case direction.Horizontal:
							x, y = gtx.Metric.PxToDp(childOffset), 0
						}

						return
					})
				}
			}()

			visibleOffset := -fw.list.Position.Offset

			return fw.list.Layout(gtx, len(children), func(gtx layout.Context, index int) layout.Dimensions {
				c := children[index]

				dims := c.Layout(gtx)

				if isVisible(index) {
					// when visible
					childrenOffsets[index] = visibleOffset
					visibleOffset += dims.Size.Y
				}

				return dims
			})
		}))
	}

	return fw.container.Layout(gtx, fw, WidgetPainterFunc(func(gtx layout.Context) layout.Dimensions {
		children := make([]*flexChild, 0)

		idx := 0
		addFlexChild := func(child *flexChild) {
			if idx > 0 && fw.Spacing != 0 {
				children = append(children, sized(Graph(func(gtx layout.Context) layout.Dimensions {
					if fw.Axis == direction.Vertical {
						return layout.Dimensions{
							Size: image.Point{
								Y: gtx.Dp(fw.Spacing),
							},
						}
					}
					return layout.Dimensions{
						Size: image.Point{
							X: gtx.Dp(fw.Spacing),
						},
					}
				})))
			}

			children = append(children, child)
			idx++
		}

		calculator := &flexCalculator{axis: fw.Axis}
		calculator.reset(gtx.Constraints.Min)

		for child := range node.IterChildElement(context.Background(), fw) {
			if w, ok := child.(Widget); ok {
				weight := float32(0)

				if getter, ok := w.(layout.WeightGetter); ok {
					if v, ok := getter.Weight(); ok {
						weight = v
					}
				}

				// force weight when arrangement is EqualWeight
				if fw.Arrangement == arrangement.EqualWeight {
					weight = 1
				}

				if weight != 0 {
					addFlexChild(flexed(weight, w))
					calculator.incrWeight(weight)
				} else {
					addFlexChild(sized(w))
				}
			}
		}

		if len(children) == 0 {
			return layout.Dimensions{}
		}

		for i := range children {
			child := children[i]

			if !child.flexed() {
				child.Paint(gtx)
				calculator.incr(child.dims.Size)
			}
		}

		remainSpacing := calculator.remainSpacing(gtx.Constraints.Max)

		for i := range children {
			child := children[i]

			if child.flexed() {
				if remainSpacing > 0 {
					gtx.Constraints = calculator.constraints(gtx.Constraints, child.weight, remainSpacing)
				}

				if sc, ok := child.widget.(paint.SizeSetter); ok {
					switch fw.Axis {
					case direction.Horizontal:
						sc.SetSize(-1, size.Width)
					case direction.Vertical:
						sc.SetSize(-1, size.Height)
					}
				}

				child.Paint(gtx)
				calculator.incr(child.dims.Size)
			}
		}

		offset := 0

		for i := range children {
			child := children[i]

			positionChild(fw, child.widget, func() (x, y unit.Dp) {
				off := image.Point{}
				off = off.Add(image.Pt(fw.offsetOfAlignment(calculator.size, child.dims.Size)))
				off = off.Add(image.Pt(fw.offsetOfArrangement(remainSpacing, len(children), i)))

				switch fw.Axis {
				case direction.Horizontal:
					off.X += offset
					offset += child.dims.Size.X
				case direction.Vertical:
					off.Y += offset
					offset += child.dims.Size.Y
				}

				defer op.Offset(image.Pt(off.X, off.Y)).Push(gtx.Ops).Pop()
				child.call.Add(gtx.Ops)

				return gtx.Metric.PxToDp(off.X), gtx.Metric.PxToDp(off.Y)
			})
		}

		return layout.Dimensions{
			Size: calculator.size,
		}
	}))
}

type flexCalculator struct {
	axis        direction.Axis
	size        image.Point
	totalWeight float32
}

func (c *flexCalculator) reset(size image.Point) {
	switch c.axis {
	case direction.Horizontal:
		c.size.Y = size.Y
	case direction.Vertical:
		c.size.X = size.X
	}
}

func (c *flexCalculator) incr(size image.Point) {
	switch c.axis {
	case direction.Horizontal:
		c.size.X += size.X

		if size.Y > c.size.Y {
			c.size.Y = size.Y
		}
	case direction.Vertical:
		c.size.Y += size.Y

		if size.X > c.size.X {
			c.size.X = size.X
		}
	}
}

func (c *flexCalculator) remainSpacing(max image.Point) int {
	switch c.axis {
	case direction.Horizontal:
		return max.X - c.size.X
	case direction.Vertical:
		return max.Y - c.size.Y
	}
	return 0
}

func (c *flexCalculator) constraints(constraints layout.Constraints, weight float32, spacing int) layout.Constraints {
	if c.totalWeight != 0 {
		switch c.axis {
		case direction.Horizontal:
			constraints.Min.X = int(math.Round(float64(float32(spacing) * weight / c.totalWeight)))
			constraints.Max.X = constraints.Min.X
		case direction.Vertical:
			constraints.Min.Y = int(math.Round(float64(float32(spacing) * weight / c.totalWeight)))
			constraints.Max.Y = constraints.Min.Y
		}
	}

	return constraints
}

func (c *flexCalculator) incrWeight(weight float32) {
	c.totalWeight += weight
}

func (fw *flexWidget) offsetOfAlignment(parent image.Point, child image.Point) (int, int) {
	switch fw.Alignment {
	case alignment.End:
		if fw.Axis == direction.Horizontal {
			return 0, parent.Y - child.Y
		}
		return parent.X - child.X, 0
	case alignment.Center:
		if fw.Axis == direction.Horizontal {
			return 0, (parent.Y - child.Y) / 2
		}
		return (parent.X - child.X) / 2, 0
	}
	return 0, 0
}

func (fw *flexWidget) offsetOfArrangement(spacing int, n int, idx int) (int, int) {
	if spacing <= 0 {
		return 0, 0
	}

	switch fw.Arrangement {
	case arrangement.SpaceEvenly:
		if fw.Axis == direction.Horizontal {
			return spacing / (n + 1) * (idx + 1), 0
		}
		return 0, spacing / (n + 1) * (idx + 1)
	case arrangement.SpaceAround:
		if fw.Axis == direction.Horizontal {
			return spacing / (n * 2) * (idx*2 + 1), 0
		}
		return 0, spacing / (n * 2) * (idx*2 + 1)
	case arrangement.SpaceBetween:
		if (n - 1) == 0 {
			return 0, 0
		}
		if fw.Axis == direction.Horizontal {
			return spacing / (n - 1) * idx, 0
		}
		return 0, spacing / (n - 1) * idx
	case arrangement.Start:
		if fw.Axis == direction.Horizontal {
			return 0, 0
		}
		return 0, 0
	case arrangement.End:
		if fw.Axis == direction.Horizontal {
			return spacing, 0
		}
		return 0, spacing
	case arrangement.Center:
		if fw.Axis == direction.Horizontal {
			return spacing / 2, 0
		}
		return 0, spacing / 2
	}
	return 0, 0
}

func sized(w Widget) *flexChild {
	return &flexChild{widget: w}
}

func flexed(weight float32, w Widget) *flexChild {
	return &flexChild{weight: weight, widget: w}
}

type flexChild struct {
	weight float32
	widget Widget

	call    op.CallOp
	dims    layout.Dimensions
	painted bool
}

func (c *flexChild) flexed() bool {
	return c.weight > 0
}

func (c *flexChild) Paint(gtx layout.Context) {
	if c.painted {
		return
	}

	c.painted = true
	c.call = paint.Group(gtx.Ops, func() {
		c.dims = c.widget.Layout(gtx)
	})
}
