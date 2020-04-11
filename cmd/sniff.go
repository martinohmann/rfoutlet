package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
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
			return options.Run()
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

func (o *SniffOptions) Run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}
	defer chip.Close()

	receiver, err := gpio.NewReceiver(chip, int(o.Pin))
	if err != nil {
		return err
	}
	defer receiver.Close()

	for {
		select {
		case res := <-receiver.Receive():
			fmt.Printf("received code=%d pulseLength=%d bitLength=%d protocol=%d\n",
				res.Code, res.PulseLength, res.BitLength, res.Protocol)
		case <-interrupt:
			fmt.Println("received interrupt")
			return nil
		}
	}
}
