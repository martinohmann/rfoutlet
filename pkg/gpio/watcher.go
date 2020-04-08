package gpio

import "github.com/warthog618/gpiod"

type watcher struct {
	pin    Closer
	events chan gpiod.LineEvent
}

// NewWatcher creates a new Watcher for the given pin offset of chip.
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

// handleEvent is the event handler that is invoked by gpiod whether the signal
// edge changes.
func (w *watcher) handleEvent(event gpiod.LineEvent) {
	w.events <- event
}

// Watch implements Watcher.
func (w *watcher) Watch() <-chan gpiod.LineEvent {
	return w.events
}

// Close implements Closer.
func (w *watcher) Close() error {
	defer close(w.events)
	return w.pin.Close()
}
