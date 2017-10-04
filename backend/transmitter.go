package backend

import (
	"fmt"
)

const gpioPin = 0

const transmitCommand = "codesend"

// Transmit transmits the given code via the configured gpio pin
func Transmit(code int, pulseLength int) error {
	fmt.Printf("transmitting code %d with pulse length %d\n", code, pulseLength)

	return nil

	// args := []string{
	// 	fmt.Sprintf("%d", code),
	// 	"-p",
	// 	fmt.Sprintf("%d", gpioPin),
	// 	"-l",
	// 	fmt.Sprintf("%d", pulseLength),
	// }

	// return exec.Command(transmitCommand, args...).Run()
}
