package gpio

import (
	"github.com/warthog618/gpiod"
)

// Pin defines an interface for a gpio pin
type Pin interface {
	SetValue(value int) error
	Reconfigure(options ...gpiod.LineConfig) error
	Close() error
}

type FakePin struct {
	Value int
	Err   error
}

func (p *FakePin) SetValue(value int) error {
	p.Value = value
	return p.Err
}

func (p *FakePin) Reconfigure(options ...gpiod.LineConfig) error {
	return p.Err
}

func (p *FakePin) Close() error {
	return p.Err
}
