// Package controller provides a controller that handles all outlet commands
// sequentially.
package controller

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "controller")

// Broadcaster can broadcast messages to all connected clients.
type Broadcaster interface {
	// Broadcast broadcasts msg to all connected clients.
	Broadcast(msg []byte)
}

// Controller controls the outlets registered to the registry.
type Controller struct {
	// Registry contains all known outlets and outlet groups.
	Registry *outlet.Registry
	// Switcher switches outlets on or off based on commands from the
	// CommandQueue.
	Switcher outlet.Switcher
	// Broadcaster broadcasts state updates to all connected clients.
	Broadcaster Broadcaster
	// CommandQueue is consumed sequentially by the controller. The commands
	// may cause outlet and group state changes which are communicated back to
	// one or more connected clients.
	CommandQueue <-chan command.Command
}

// Run runs the main control loop until stopCh is closed.
func (c *Controller) Run(stopCh <-chan struct{}) {
	for {
		select {
		case cmd, ok := <-c.CommandQueue:
			if !ok {
				log.Error("command queue was closed unexpectedly, shutting down controller")
				return
			}

			err := c.handleCommand(cmd)
			if err != nil {
				log.WithField("command", fmt.Sprintf("%T", cmd)).
					Errorf("error handling command: %v", err)
			}
		case <-stopCh:
			log.Info("shutting down controller")
			return
		}
	}
}

// commandContext creates a new command.Context.
func (c *Controller) commandContext() command.Context {
	return command.Context{
		Registry: c.Registry,
		Switcher: c.Switcher,
	}
}

// handleCommand executes cmd and may trigger broadcasts of state changes back
// to the connected clients.
func (c *Controller) handleCommand(cmd command.Command) error {
	log.WithField("command", fmt.Sprintf("%T", cmd)).
		Debug("handling command")

	ctx := c.commandContext()

	broadcast, err := cmd.Execute(ctx)
	if err != nil || !broadcast {
		return err
	}

	return c.broadcastState()
}

// broadcastState broadcasts the current outlet group state back to connected
// clients.
func (c *Controller) broadcastState() error {
	msg, err := json.Marshal(c.Registry.GetGroups())
	if err != nil {
		return err
	}

	c.Broadcaster.Broadcast(msg)

	return nil
}
