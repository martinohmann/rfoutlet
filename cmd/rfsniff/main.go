package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brian-armstrong/gpio"
)

var (
	pulseLength = flag.Int("pulse-length", 0, "pulse length")
	gpioPin     = flag.Int("gpio-pin", 27, "gpio pin to sniff on")
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

	watcher := gpio.NewWatcher()
	watcher.AddPin(uint(*gpioPin))
	defer watcher.Close()

	for {
		pin, value := watcher.Watch()
		fmt.Printf("read %d from gpio %d\n", value, pin)
	}
}
