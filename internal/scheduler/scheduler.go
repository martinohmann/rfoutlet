package scheduler

import (
	"log"
	"time"

	"github.com/martinohmann/rfoutlet/internal/outlet"
)

// Scheduler type definition
type Scheduler struct {
	outlets  map[*outlet.Outlet]bool
	switcher outlet.Switcher
	ticker   *time.Ticker
	outlet   chan *outlet.Outlet
}

// New creates a new scheduler
func New(switcher outlet.Switcher) *Scheduler {
	return NewWithInterval(switcher, time.Second)
}

// NewWithInterval create a new scheduler with user defined check interval
func NewWithInterval(switcher outlet.Switcher, interval time.Duration) *Scheduler {
	s := &Scheduler{
		outlets:  make(map[*outlet.Outlet]bool),
		switcher: switcher,
		ticker:   time.NewTicker(interval),
		outlet:   make(chan *outlet.Outlet),
	}

	go s.run()

	return s
}

// Register registers an outlet to the scheduler
func (s *Scheduler) Register(outlet *outlet.Outlet) {
	s.outlet <- outlet
}

func (s *Scheduler) run() {
	for {
		select {
		case outlet := <-s.outlet:
			s.outlets[outlet] = true
		case <-s.ticker.C:
			s.check()
		}
	}
}

func (s *Scheduler) check() {
	for o := range s.outlets {
		if o.Schedule == nil || !o.Schedule.Enabled() {
			continue
		}

		if o.Schedule.Contains(time.Now()) {
			s.applyState(o, outlet.StateOn)
		} else {
			s.applyState(o, outlet.StateOff)
		}
	}
}

func (s *Scheduler) applyState(o *outlet.Outlet, newState outlet.State) {
	if o.GetState() == newState {
		return
	}

	if err := s.switcher.Switch(o, newState); err != nil {
		log.Println(err)
	}
}
