package gpio

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "gpio: ", log.LstdFlags|log.Lshortfile)
}

type CodeTransmitter interface {
	Transmit(uint64, int) error
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
func (t *CodesendTransmitter) Transmit(code uint64, pulseLength int) error {
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

type NativeTransmitter struct {
	gpioPin int
}

func NewNativeTransmitter(gpioPin int) *NativeTransmitter {
	return &NativeTransmitter{
		gpioPin: gpioPin,
	}
}

// Transmit transmits the given code via the configured gpio pin
func (t *NativeTransmitter) Transmit(code uint64, pulseLength int) error {
	logger.Printf("transmitting code=%d pulseLength=%d\n", code, pulseLength)

	err := rpio.Open()
	if err != nil {
		return err
	}

	defer rpio.Close()

	pin := rpio.Pin(t.gpioPin)
	pin.Output()

	for i := 0; i < 10; i++ {
		for j := 24 - 1; j >= 0; j-- {
			if code&(1<<uint64(j)) > 0 {
				t.send1(pin, pulseLength)
			} else {
				t.send0(pin, pulseLength)
			}
		}
		t.sendSync(pin, pulseLength)
	}

	return nil
}

func (t *NativeTransmitter) send0(pin rpio.Pin, pulseLength int) {
	t.transmit(pin, 1, 3, pulseLength)
}

func (t *NativeTransmitter) send1(pin rpio.Pin, pulseLength int) {
	t.transmit(pin, 3, 1, pulseLength)
}

func (t *NativeTransmitter) sendSync(pin rpio.Pin, pulseLength int) {
	t.transmit(pin, 1, 31, pulseLength)
}

func (t *NativeTransmitter) transmit(pin rpio.Pin, highPulses int, lowPulses int, pulseLength int) {
	pin.High()
	time.Sleep(time.Microsecond * time.Duration(pulseLength*highPulses))
	pin.Low()
	time.Sleep(time.Microsecond * time.Duration(pulseLength*lowPulses))
}
