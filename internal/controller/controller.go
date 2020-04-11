package controller

import (
	"encoding/json"
	"log"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

type Broadcaster interface {
	Broadcast(msg []byte)
}

type Controller struct {
	Registry     *outlet.Registry
	Switcher     outlet.Switcher
	Broadcaster  Broadcaster
	CommandQueue <-chan command.Command
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	for {
		select {
		case cmd := <-c.CommandQueue:
			err := c.handleCommand(cmd)
			if err != nil {
				log.Println(err)
			}
		case <-stopCh:
			return
		}
	}
}

func (c *Controller) commandContext() command.Context {
	return command.Context{
		Registry: c.Registry,
		Switcher: c.Switcher,
	}
}

func (c *Controller) handleCommand(cmd command.Command) error {
	ctx := c.commandContext()

	broadcast, err := cmd.Execute(ctx)
	if err != nil || !broadcast {
		return err
	}

	return c.broadcastState()
}

func (c *Controller) broadcastState() error {
	msg, err := json.Marshal(c.Registry.GetGroups())
	if err != nil {
		return err
	}

	c.Broadcaster.Broadcast(msg)

	return nil
}
