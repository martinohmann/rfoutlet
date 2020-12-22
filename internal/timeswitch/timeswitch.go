// Package timeswitch implements the time switch logic for outlets. Outlets can
// be configured to automatically be turned on or off during certain periods of
// the day or week. The time switch will emit state correction commands based
// on the schedule that may be defined on an outlet.
package timeswitch

import (
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "timeswitch")

// TimeSwitch checks if outlets should be enabled or disabled based on their
// schedule and send out commands to bring them to the desired state.
type TimeSwitch struct {
	Registry     *outlet.Registry
	CommandQueue chan<- command.Command
	Clock        clockwork.Clock
}

// New creates a new *TimeSwitch which will observe the outlets in the registry
// and eventually push commands to the queue if a state change is required.
func New(registry *outlet.Registry, queue chan<- command.Command) *TimeSwitch {
	return &TimeSwitch{
		Registry:     registry,
		CommandQueue: queue,
		Clock:        clockwork.NewRealClock(),
	}
}

// Run runs the time switch control loop which periodically checks if an outlet
// should be enabled or disabled. Whenever an outlet should change its state,
// it will push a TimeSwitchCommand into the CommandQueue.
func (s *TimeSwitch) Run(stopCh <-chan struct{}) {
	for {
		select {
		case <-s.Clock.After(s.secondsUntilNextMinute()):
			s.check()
		case <-stopCh:
			log.Info("shutting down time switch")
			return
		}
	}
}

// secondsUntilNextMinute returns the seconds until the next minute starts.
func (s *TimeSwitch) secondsUntilNextMinute() time.Duration {
	return time.Second * time.Duration(60-s.Clock.Now().Second())
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
			s.CommandQueue <- command.StateCorrectionCommand{
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
