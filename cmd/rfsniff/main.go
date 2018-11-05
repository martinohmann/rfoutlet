package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

var (
	gpioPin = flag.Uint("gpio-pin", gpio.DefaultReceivePin, "gpio pin to sniff on")
	usage   = func() {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func init() {
	flag.Usage = usage
}

func main() {
	flag.Parse()

	receiver := gpio.NewReceiver(*gpioPin)
	defer receiver.Close()

	for res := range receiver.Receive() {
		fmt.Printf("received code=%d pulseLength=%d bitLength=%d protocol=%d\n",
			res.Code, res.PulseLength, res.BitLength, res.Protocol)
	}
}
