// Package statedrift provides a detector for listening to rf codes sent out by
// anything else than rfoutlet (e.g. the physical remote control for the
// outlet).
package statedrift

import (
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/controller/commands"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "statedrift")

// Detector sniffs for sent out rf codes using the receiver and pushes state
// correction commands into the command queue if necessary. This allows for the
// detection of codes sent out by pressing buttons on a physical outlet remote
// control.
type Detector struct {
	Registry     *outlet.Registry
	Receiver     gpio.CodeReceiver
	CommandQueue chan<- command.Command
}

// NewDetector creates a new *Detector.
func NewDetector(registry *outlet.Registry, receiver gpio.CodeReceiver, queue chan<- command.Command) *Detector {
	return &Detector{
		Registry:     registry,
		Receiver:     receiver,
		CommandQueue: queue,
	}
}

// Run runs the state drift detection loop until stopCh is closed.
func (d *Detector) Run(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			log.Info("shutting down state drift detector")
			return
		case result, ok := <-d.Receiver.Receive():
			if !ok {
				log.Error("receiver was closed unexpectedly, shutting down state drift detector")
				return
			}

			var found bool

			for _, o := range d.Registry.GetOutlets() {
				if result.Code == o.CodeOn && o.GetState() != outlet.StateOn {
					found = true
					d.CommandQueue <- commands.StateCorrectionCommand{
						Outlet:       o,
						DesiredState: outlet.StateOn,
					}
				} else if result.Code == o.CodeOff && o.GetState() != outlet.StateOff {
					found = true
					d.CommandQueue <- commands.StateCorrectionCommand{
						Outlet:       o,
						DesiredState: outlet.StateOff,
					}
				}

				if found {
					break
				}
			}
		}
	}
}
