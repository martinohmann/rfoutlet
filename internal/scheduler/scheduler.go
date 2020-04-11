package scheduler

import (
	"time"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

type Scheduler struct {
	Interval     time.Duration
	Registry     *outlet.Registry
	CommandQueue chan<- command.Command
}

func New(registry *outlet.Registry, queue chan<- command.Command) *Scheduler {
	return &Scheduler{
		Interval:     10 * time.Second,
		Registry:     registry,
		CommandQueue: queue,
	}
}

func (s *Scheduler) Run(stopCh <-chan struct{}) {
	ticker := time.NewTicker(s.Interval)

	for {
		select {
		case <-ticker.C:
			s.schedule()
		case <-stopCh:
			ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) schedule() {
	for _, outlet := range s.Registry.GetOutlets() {
		if !outlet.Schedule.Enabled() {
			continue
		}

		desiredState := getDesiredState(outlet)

		if outlet.GetState() != desiredState {
			s.CommandQueue <- command.ScheduleCommand{
				Outlet:       outlet,
				DesiredState: desiredState,
			}
		}
	}
}

func getDesiredState(o *outlet.Outlet) outlet.State {
	if o.Schedule.Contains(time.Now()) {
		return outlet.StateOn
	}

	return outlet.StateOff
}
