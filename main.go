package main

import (
	"fmt"
	"os"

	"github.com/martinohmann/rfoutlet/cmd"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "rfoutlet",
		Short:         "A tool for interacting with remote controlled outlets",
		Long:          "rfoutlet is a tool for interacting with remote controlled outlets. It provides functionality to sniff and transmit the codes controlling the outlets.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	return cmd
}

func main() {
	rootCmd := newRootCommand()

	rootCmd.AddCommand(cmd.NewServeCommand())
	rootCmd.AddCommand(cmd.NewSniffCommand())
	rootCmd.AddCommand(cmd.NewTransmitCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
