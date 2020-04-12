package cmd

import (
	"fmt"
	"strconv"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
)

func NewTransmitCommand() *cobra.Command {
	options := &TransmitOptions{
		PulseLength: config.DefaultPulseLength,
		Pin:         config.DefaultTransmitPin,
		Protocol:    config.DefaultProtocol,
	}

	cmd := &cobra.Command{
		Use:   "transmit [codes...]",
		Short: "Send out codes to remote controlled outlets",
		Long:  "The transmit command can be used send out codes to remote controlled outlets.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.Run(args)
		},
	}

	options.AddFlags(cmd)

	return cmd
}

type TransmitOptions struct {
	config.Config
	PulseLength uint
	Pin         uint
	Protocol    int
}

func (o *TransmitOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().UintVar(&o.PulseLength, "pulse-length", o.PulseLength, "pulse length")
	cmd.Flags().UintVar(&o.Pin, "pin", o.Pin, "gpio pin to transmit on")
	cmd.Flags().IntVar(&o.Protocol, "protocol", o.Protocol, "transmission protocl")
}

func (o *TransmitOptions) Run(args []string) error {
	if o.Protocol < 1 || o.Protocol > len(gpio.DefaultProtocols) {
		return fmt.Errorf("protocol %d does not exist", o.Protocol)
	}

	proto := gpio.DefaultProtocols[o.Protocol-1]

	codes := make([]uint64, len(args))

	for i, arg := range args {
		var err error
		codes[i], err = strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse code: %v", err)
		}
	}

	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return fmt.Errorf("failed to open gpio device: %v", err)
	}
	defer chip.Close()

	transmitter, err := gpio.NewTransmitter(chip, int(o.Pin))
	if err != nil {
		return fmt.Errorf("failed to create gpio transmitter: %v", err)
	}
	defer transmitter.Close()

	stopCh := make(chan struct{})

	go handleSignals(stopCh)

	log := log.WithFields(log.Fields{
		"pulseLength": o.PulseLength,
		"protocol":    o.Protocol,
	})

	for _, code := range codes {
		log.Infof("transmitting code %d", code)

		select {
		case <-transmitter.Transmit(code, proto, o.PulseLength):
		case <-stopCh:
			return nil
		}
	}

	return nil
}
