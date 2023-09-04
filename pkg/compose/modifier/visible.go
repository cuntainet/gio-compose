package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
)

func Visible(v bool) modifier.Modifier[any] {
	return &visibleModifier{visible: v}
}

type visibleModifier struct {
	visible bool
}

func (v *visibleModifier) Modify(target any) {
	if s, ok := target.(paint.VisibleSetter); ok {
		s.SetVisible(v.visible)
	}
}
