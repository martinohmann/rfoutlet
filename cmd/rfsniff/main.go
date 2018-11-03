package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

var (
	pulseLength = flag.Int("pulse-length", 0, "pulse length")
	gpioPin     = flag.Uint("gpio-pin", gpio.DefaultReceivePin, "gpio pin to sniff on")
	usage       = func() {
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

	receiver.Receive(func(code uint64, pulseLength int64, bitLength uint, protocol int) {
		fmt.Printf("received code=%d pulseLength=%d bitLength=%d protocol=%d\n", code, pulseLength, bitLength, protocol)
	})

	receiver.Wait()
}
