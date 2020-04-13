// rfoutlet provides outlet control via cli and web interface for
// Raspberry PI 2/3.
//
// The transmitter and receiver logic has been ported from the great
// https://github.com/sui77/rc-switch C++ project to golang.
//
// rfoutlet comes with ready to use commands for transmitting and receiving
// remote control codes as well as a command for serving a web interface (see
// cmd/ directory). The pkg/ directory exposes the gpio package which contains
// the receiver and transmitter code.
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "rfoutlet",
		Short:         "A tool for interacting with remote controlled outlets",
		Long:          "rfoutlet is a tool for interacting with remote controlled outlets. It provides functionality to sniff and transmit the codes controlling the outlets.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			gin.SetMode(gin.ReleaseMode)

			if debug {
				gin.SetMode(gin.DebugMode)
				log.SetLevel(log.DebugLevel)
				log.SetFormatter(&log.TextFormatter{
					FullTimestamp: true,
				})
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&debug, "debug", debug, "enable debug mode. this will cause more verbose output")

	return cmd
}

func main() {
	rootCmd := newRootCommand()

	rootCmd.AddCommand(cmd.NewServeCommand())
	rootCmd.AddCommand(cmd.NewSniffCommand())
	rootCmd.AddCommand(cmd.NewTransmitCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
