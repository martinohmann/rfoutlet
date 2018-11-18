package scheduler

import (
	"time"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/state"
)

// Scheduler type definition
type Scheduler struct {
	ctx     *context.Context
	control *control.Control
	ticker  *time.Ticker
	stop    chan bool
}

// New creates a new scheduler
func New(ctx *context.Context, control *control.Control, interval time.Duration) *Scheduler {
	return &Scheduler{
		ctx:     ctx,
		control: control,
		ticker:  time.NewTicker(interval),
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
	for _, g := range s.ctx.Groups {
		for _, o := range g.Outlets {
			if o.Schedule == nil {
				continue
			}

			if o.Schedule.Contains(time.Now()) {
				s.maybeSwitchState(o, state.SwitchStateOn)
			} else {
				s.maybeSwitchState(o, state.SwitchStateOff)
			}
		}
	}
}

func (s *Scheduler) maybeSwitchState(o *context.Outlet, newState state.SwitchState) error {
	if o.State == newState {
		return nil
	}

	return s.control.SwitchState(o, newState)
}
