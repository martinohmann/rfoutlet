package control

import (
	"time"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/state"
)

// Scheduler type definition
type Scheduler struct {
	ctx    *context.Context
	ticker *time.Ticker
	stop   chan bool
}

// NewScheduler create a new scheduler
func NewScheduler(ctx *context.Context, interval time.Duration) *Scheduler {
	return &Scheduler{
		ctx:    ctx,
		ticker: time.NewTicker(interval),
		stop:   make(chan bool, 1),
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
				if o.State != state.SwitchStateOn {
					// switch on
				}
			} else if o.State != state.SwitchStateOff {
				// switch off
			}
		}
	}
}
