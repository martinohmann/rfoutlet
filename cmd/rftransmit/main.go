package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

var (
	pulseLength = flag.Uint("pulse-length", gpio.DefaultPulseLength, "pulse length")
	gpioPin     = flag.Uint("gpio-pin", gpio.DefaultTransmitPin, "gpio pin to transmit on")
	protocol    = flag.Int("protocol", gpio.DefaultProtocol, "transmission protocol")
	usage       = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <code>\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func init() {
	flag.Usage = usage
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	c, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	code := uint64(c)

	t, err := gpio.NewTransmitter(*gpioPin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer t.Close()

	if err = t.Transmit(code, *protocol, *pulseLength); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	<-t.Transmitted()
}
