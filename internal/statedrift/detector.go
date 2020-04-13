package statedrift

import (
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "statedrift")

type Detector struct {
	Registry     *outlet.Registry
	Receiver     gpio.CodeReceiver
	CommandQueue chan<- command.Command
}

func NewDetector(registry *outlet.Registry, receiver gpio.CodeReceiver, queue chan<- command.Command) *Detector {
	return &Detector{
		Registry:     registry,
		Receiver:     receiver,
		CommandQueue: queue,
	}
}

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

			for _, o := range d.Registry.GetOutlets() {
				if result.Code == o.CodeOn && o.GetState() != outlet.StateOn {
					d.CommandQueue <- StateCorrectionCommand{
						Outlet:       o,
						DesiredState: outlet.StateOn,
					}
				} else if result.Code == o.CodeOff && o.GetState() != outlet.StateOff {
					d.CommandQueue <- StateCorrectionCommand{
						Outlet:       o,
						DesiredState: outlet.StateOff,
					}
				}
			}
		}
	}
}
