// Package command provides the interfaces for commands that can be pushed
// into a command queue for the controller to consume.
package command

import "github.com/martinohmann/rfoutlet/internal/outlet"

// Context is passed to every command.
type Context struct {
	// Registry contains all known outlets and outlet groups.
	*outlet.Registry
	// Switcher can switch an outlet on or off.
	outlet.Switcher
}

// Command is something that can be put into a command queue and is executed by
// a controller sequentially.
type Command interface {
	// Execute executes the command. The passed in context contains the outlet
	// registry and switcher to be able to manipulate the state of outlets.
	// Must return true if a potential state change should be broadcasted to
	// all connected clients.
	Execute(context Context) (broadcast bool, err error)
}

// Sender can send messages.
type Sender interface {
	// Send sends out a message.
	Send(msg []byte)
}

// SenderAwareCommand is aware of the comand's sender.
type SenderAwareCommand interface {
	Command

	// SetSender sets the sender on the command. The sender can be used to send
	// messages back to the client that issued the command.
	SetSender(sender Sender)
}

// NewTestContext creates a new Context which can be used in tests. It returns
// the wrapped registry and switcher as 2nd and 3rd return value.
func NewTestContext() (Context, *outlet.Registry, *outlet.FakeSwitch) {
	r := outlet.NewRegistry()
	s := &outlet.FakeSwitch{}

	return Context{Registry: r, Switcher: s}, r, s
}
