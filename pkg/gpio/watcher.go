package gpio

import "github.com/warthog618/gpiod"

// Watcher defines the interface for a gpio pin watcher
type Watcher interface {
	Watch() <-chan gpiod.LineEvent
	Close() error
}

type watcher struct {
	pin    *gpiod.Line
	events chan gpiod.LineEvent
}

func NewWatcher(chip *gpiod.Chip, offset int) (Watcher, error) {
	w := &watcher{
		events: make(chan gpiod.LineEvent),
	}

	pin, err := chip.RequestLine(offset, gpiod.WithBothEdges(w.handleEvent))
	if err != nil {
		return nil, err
	}

	w.pin = pin

	return w, nil
}

func (w *watcher) handleEvent(event gpiod.LineEvent) {
	select {
	case w.events <- event:
	default: // don't block if there is no consumer
	}
}

func (w *watcher) Watch() <-chan gpiod.LineEvent {
	return w.events
}

func (w *watcher) Close() error {
	defer close(w.events)
	return w.pin.Close()
}
