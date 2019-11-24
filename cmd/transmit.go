package cmd

import (
	"fmt"
	"strconv"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/spf13/cobra"
)

func NewTransmitCommand() *cobra.Command {
	options := &TransmitOptions{
		PulseLength: gpio.DefaultPulseLength,
		GpioPin:     gpio.DefaultReceivePin,
		Protocol:    gpio.DefaultProtocol,
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
	c, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	code := uint64(c)

	t := gpio.NewTransmitter(o.GpioPin)
	defer t.Close()

	fmt.Printf("transmitting code=%d pulseLength=%d protocol=%d\n", code, o.PulseLength, o.Protocol)

	err = t.Transmit(code, o.Protocol, o.PulseLength)
	if err != nil {
		return err
	}

	t.Wait()

	return nil
}
