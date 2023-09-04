package textinput

type ChangeEventProvider interface {
	WatchChangeEvent(func(text string))
	TriggerChange(text string)
}

type Events struct {
	changeEventConsumers []func(text string)
}

var _ ChangeEventProvider = &Events{}

func (events *Events) WatchChangeEvent(action func(text string)) {
	events.changeEventConsumers = append(events.changeEventConsumers, action)
}

func (events *Events) TriggerChange(text string) {
	for _, c := range events.changeEventConsumers {
		if c != nil {
			c(text)
		}
	}
}
