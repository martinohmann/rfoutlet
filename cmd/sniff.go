package cmd

import (
	"context"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewSniffCommand() *cobra.Command {
	options := &SniffOptions{
		Pin: config.DefaultReceivePin,
	}

	cmd := &cobra.Command{
		Use:   "sniff",
		Short: "Sniff codes sent out to remote controlled outlets",
		Long:  "The sniff command can be used to sniff codes sent out to remote controlled outlets.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return options.Run(cmd)
		},
	}

	options.AddFlags(cmd)

	return cmd
}

type SniffOptions struct {
	Pin uint
}

func (o *SniffOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().UintVar(&o.Pin, "pin", o.Pin, "gpio pin to sniff on")
}

func (o *SniffOptions) Run(cmd *cobra.Command) error {
	device, err := openGPIODevice(cmd)
	if err != nil {
		return err
	}
	defer device.Close()

	receiver, err := gpio.NewReceiver(device.Chip, int(o.Pin))
	if err != nil {
		return fmt.Errorf("failed to create gpio receiver: %v", err)
	}
	defer receiver.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSignals(cancel)

	for {
		select {
		case res := <-receiver.Receive():
			log.WithFields(log.Fields{
				"pulseLength": res.PulseLength,
				"protocol":    res.Protocol,
				"bitlength":   res.BitLength,
			}).Infof("received code %d", res.Code)
		case <-ctx.Done():
			return nil
		}
	}
}
