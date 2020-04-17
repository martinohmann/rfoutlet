package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
		Count:       gpio.DefaultTransmissionCount,
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
	Count       int
	Delay       time.Duration
	Infinite    bool
}

func (o *TransmitOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().UintVar(&o.PulseLength, "pulse-length", o.PulseLength, "pulse length")
	cmd.Flags().UintVar(&o.Pin, "pin", o.Pin, "gpio pin to transmit on")
	cmd.Flags().IntVar(&o.Protocol, "protocol", o.Protocol, "protocol to use for the transmission")
	cmd.Flags().IntVar(&o.Count, "count", o.Count, "number of times a code should be transmitted in a row. The higher the value, the more likely it is that an outlet actually received the code")
	cmd.Flags().DurationVar(&o.Delay, "delay", o.Delay, "delay between code transmissions")
	cmd.Flags().BoolVar(&o.Infinite, "infinite", o.Infinite, "restart the transmission of codes after the last one was sent")
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

	chip, err := gpiod.NewChip(gpioChipName)
	if err != nil {
		return fmt.Errorf("failed to open gpio device: %v", err)
	}
	defer chip.Close()

	transmitter, err := gpio.NewTransmitter(chip, int(o.Pin), gpio.TransmissionCount(o.Count))
	if err != nil {
		return fmt.Errorf("failed to create gpio transmitter: %v", err)
	}
	defer transmitter.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSignals(cancel)

	log.WithFields(log.Fields{
		"pulseLength": o.PulseLength,
		"protocol":    o.Protocol,
		"delay":       o.Delay,
		"count":       o.Count,
	}).Infof("starting transmission")

Loop:
	for _, code := range codes {
		log.Infof("transmitting code %d", code)

		select {
		case <-transmitter.Transmit(code, proto, o.PulseLength):
			<-time.After(o.Delay)
		case <-ctx.Done():
			return nil
		}
	}

	if o.Infinite {
		goto Loop
	}

	return nil
}
