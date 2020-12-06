package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/mockup"
)

type device struct {
	*gpiod.Chip
	*mockup.Mockup
}

func (d *device) Close() error {
	err := d.Chip.Close()
	if err != nil {
		return err
	}

	if d.Mockup != nil {
		log.Debug("removing gpio mockup")
		return d.Mockup.Close()
	}

	return nil
}

func openGPIODevice(cmd *cobra.Command) (*device, error) {
	gpioChipName, _ := cmd.Flags().GetString("gpio-chip")
	gpioMockup, _ := cmd.Flags().GetBool("gpio-mockup")

	var (
		dev *device = &device{}
		err error
	)

	if gpioMockup {
		log.Debug("creating gpio mockup")
		dev.Mockup, err = mockup.New([]int{40}, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create gpio mockup: %v", err)
		}
	}

	dev.Chip, err = gpiod.NewChip(gpioChipName)
	if err != nil {
		return nil, fmt.Errorf("failed to open gpio device: %v", err)
	}

	return dev, nil
}
