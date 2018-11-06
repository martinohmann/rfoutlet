// The rfsniff command can be used to sniff codes sent out by remotes for
// remote controlled outlet. Start the command and press the buttons on
// the remote. You should see the received code, pulse length, bit length
// and remote protocol in the terminal.
//
// Available command line flags:
//
//  -gpio-pin uint
//        gpio pin to sniff on (default 27)
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
