package cmd

import (
	"fmt"
	"strconv"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
)

func NewTransmitCommand() *cobra.Command {
	options := &TransmitOptions{
		PulseLength: config.DefaultPulseLength,
		GpioPin:     config.DefaultReceivePin,
		Protocol:    config.DefaultProtocol,
	}

	cmd := &cobra.Command{
		Use:   "transmit <code>",
		Short: "Send out codes to remote controlled outlets",
		Long:  "The transmit command can be used send out codes to remote controlled outlets.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.Run(args)
		},
	}

	options.AddFlags(cmd)

	return cmd
}

type TransmitOptions struct {
	PulseLength uint
	GpioPin     uint
	Protocol    int
}

func (o *TransmitOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().UintVar(&o.PulseLength, "pulse-length", o.PulseLength, "pulse length")
	cmd.Flags().UintVar(&o.GpioPin, "gpio-pin", o.GpioPin, "gpio pin to transmit on")
	cmd.Flags().IntVar(&o.Protocol, "protocol", o.Protocol, "transmission protocl")
}

func (o *TransmitOptions) Run(args []string) error {
	code, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}
	defer chip.Close()

	transmitter, err := gpio.NewTransmitter(chip, int(o.GpioPin))
	if err != nil {
		return err
	}
	defer transmitter.Close()

	fmt.Printf("transmitting code=%d pulseLength=%d protocol=%d\n", code, o.PulseLength, o.Protocol)

	err = transmitter.Transmit(code, o.Protocol, o.PulseLength)
	if err != nil {
		return err
	}

	transmitter.Wait()

	return nil
}
