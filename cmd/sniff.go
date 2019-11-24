package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/spf13/cobra"
)

func NewSniffCommand() *cobra.Command {
	options := &SniffOptions{
		GpioPin: gpio.DefaultReceivePin,
	}

	cmd := &cobra.Command{
		Use:   "sniff",
		Short: "Sniff codes sent out to remote controlled outlets",
		Long:  "The sniff command can be used to sniff codes sent out to remote controlled outlets.",
		Run: func(cmd *cobra.Command, _ []string) {
			options.Run()
		},
	}

	options.AddFlags(cmd)

	return cmd
}

type SniffOptions struct {
	GpioPin uint
}

func (o *SniffOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().UintVar(&o.GpioPin, "gpio-pin", o.GpioPin, "gpio pin to sniff on")
}

func (o *SniffOptions) Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	receiver := gpio.NewReceiver(o.GpioPin)
	defer receiver.Close()

	for {
		select {
		case res := <-receiver.Receive():
			fmt.Printf("received code=%d pulseLength=%d bitLength=%d protocol=%d\n",
				res.Code, res.PulseLength, res.BitLength, res.Protocol)
		case <-interrupt:
			fmt.Println("received interrupt")
			return
		}
	}
}
