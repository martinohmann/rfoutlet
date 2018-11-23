package scheduler

import (
	"encoding/json"
	"time"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/state"
)

// Scheduler type definition
type Scheduler struct {
	control *control.Control
	ticker  *time.Ticker
	stop    chan bool
}

// New creates a new scheduler
func New(control *control.Control) *Scheduler {
	return &Scheduler{
		control: control,
		ticker:  time.NewTicker(10 * time.Second),
		stop:    make(chan bool, 1),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	go s.run()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.stop <- true
}

func (s *Scheduler) run() {
	for {
		select {
		case <-s.ticker.C:
			s.schedule()
		case <-s.stop:
			s.ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) schedule() {
	for _, g := range s.control.Groups() {
		for _, o := range g.Outlets {
			sch := o.GetSchedule()

			if sch == nil || !sch.Enabled() {
				continue
			}

			if sch.Contains(time.Now()) {
				s.transitionToState(o, state.SwitchStateOn)
			} else {
				s.transitionToState(o, state.SwitchStateOff)
			}
		}
	}
}

func (s *Scheduler) transitionToState(o *context.Outlet, newState state.SwitchState) error {
	if o.State == newState {
		return nil
	}

	if err := s.control.SwitchState(o, newState); err != nil {
		return err
	}

	b, err := json.Marshal(s.control.Groups())
	if err != nil {
		return err
	}

	s.control.Broadcast(b)

	return nil
}
