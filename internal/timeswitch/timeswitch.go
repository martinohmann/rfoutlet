// Package timeswitch implements the time switch logic for outlets. Outlets can
// be configured to automatically be turned on or off during certain periods of
// the day or week. The time switch will emit state correction commands based
// on the schedule that may be defined on an outlet.
package timeswitch

import (
	"time"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/controller/commands"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "timeswitch")

// TimeSwitch checks if outlets should be enabled or disabled based on their
// schedule and send out commands to bring them to the desired state.
type TimeSwitch struct {
	Interval     time.Duration
	Registry     *outlet.Registry
	CommandQueue chan<- command.Command
}

// New creates a new *TimeSwitch which will observe the outlets in the registry
// and eventually push commands to the queue if a state change is required.
func New(registry *outlet.Registry, queue chan<- command.Command) *TimeSwitch {
	return &TimeSwitch{
		Interval:     10 * time.Second,
		Registry:     registry,
		CommandQueue: queue,
	}
}

// Run runs the time switch control loop which periodically checks if an outlet
// should be enabled or disabled. Whenever an outlet should change its state,
// it will push a TimeSwitchCommand into the CommandQueue.
func (s *TimeSwitch) Run(stopCh <-chan struct{}) {
	ticker := time.NewTicker(s.Interval)

	for {
		select {
		case <-ticker.C:
			s.check()
		case <-stopCh:
			ticker.Stop()
			log.Info("shutting down time switch")
			return
		}
	}
}

func (s *TimeSwitch) check() {
	for _, outlet := range s.Registry.GetOutlets() {
		if !outlet.Schedule.Enabled() {
			continue
		}

		desiredState := getDesiredState(outlet)

		// We only send out commands if the outlet is not in the desired state
		// to avoid spamming the command queue.
		if outlet.GetState() != desiredState {
			s.CommandQueue <- commands.StateCorrectionCommand{
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
