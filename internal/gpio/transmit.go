package gpio

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const transmitCommand = "codesend"

var logger *log.Logger

type TransmitFunc func(int, int, int) error

func init() {
	logger = log.New(os.Stdout, "gpio: ", log.LstdFlags|log.Lshortfile)
}

// Transmit transmits the given code via the configured gpio pin
func Transmit(code int, gpioPin int, pulseLength int) error {
	logger.Printf("transmitting code=%d pulseLength=%d gpioPin=%d\n", code, pulseLength, gpioPin)

	args := []string{
		fmt.Sprintf("%d", code),
		"-p",
		fmt.Sprintf("%d", gpioPin),
		"-l",
		fmt.Sprintf("%d", pulseLength),
	}

	return exec.Command(transmitCommand, args...).Run()
}
