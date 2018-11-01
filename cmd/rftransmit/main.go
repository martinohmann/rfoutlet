package main

import (
	"flag"
	"strconv"

	"github.com/martinohmann/rfoutlet/internal/gpio"
)

const (
	defaultPulseLength = 189
	defaultGpioPin     = 0
)

var (
	pulseLength = flag.Int("pulse-length", defaultPulseLength, "pulse length")
	gpioPin     = flag.Int("gpio-pin", defaultGpioPin, "gpio pin to transmit on")
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

	t := gpio.NewNativeTransmitter(*gpioPin)

	err = t.Transmit(code, *pulseLength)
	if err != nil {
		panic(err)
	}
}
