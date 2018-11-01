package gpio

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "gpio: ", log.LstdFlags|log.Lshortfile)
}

type CodeTransmitter interface {
	Transmit(int, int) error
}

type CodesendTransmitter struct {
	gpioPin int
}

func NewCodesendTransmitter(gpioPin int) *CodesendTransmitter {
	return &CodesendTransmitter{
		gpioPin: gpioPin,
	}
}

// Transmit transmits the given code via the configured gpio pin
func (t *CodesendTransmitter) Transmit(code int, pulseLength int) error {
	logger.Printf("transmitting code=%d pulseLength=%d\n", code, pulseLength)

	args := []string{
		fmt.Sprintf("%d", code),
		"-p",
		fmt.Sprintf("%d", t.gpioPin),
		"-l",
		fmt.Sprintf("%d", pulseLength),
	}

	return exec.Command("codesend", args...).Run()
}
