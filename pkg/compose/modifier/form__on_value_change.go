package modifier

import (
	"github.com/octohelm/gio-compose/pkg/gesture/textinput"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func OnValueChange(action func(v string)) modifier.Modifier[any] {
	return &valueChangeWatcher{
		action: action,
	}
}

type valueChangeWatcher struct {
	action func(v string)
}

func (w *valueChangeWatcher) Modify(target any) {
	if p, ok := target.(textinput.ChangeEventProvider); ok {
		p.WatchChangeEvent(w.action)
	}
}
