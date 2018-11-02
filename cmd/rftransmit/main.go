package main

import (
	"flag"
	"strconv"

	"github.com/martinohmann/rfoutlet/internal/gpio"
)

var (
	pulseLength = flag.Int("pulse-length", gpio.DefaultPulseLength, "pulse length")
	gpioPin     = flag.Int("gpio-pin", gpio.DefaultGpioPin, "gpio pin to transmit on")
	protocol    = flag.Int("protocol", gpio.DefaultProtocol, "transmission protocol")
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		panic("code missing")
	}

	c, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}

	code := uint64(c)

	t, err := gpio.NewTransmitter(*gpioPin)
	if err != nil {
		panic(err)
	}

	defer t.Close()

	err = t.Transmit(code, *protocol, *pulseLength)
	if err != nil {
		panic(err)
	}
}
