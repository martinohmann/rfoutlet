package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/imdario/mergo"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/controller"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/scheduler"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/internal/websocket"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
)

const webDir = "../web/build"

func NewServeCommand() *cobra.Command {
	options := &ServeOptions{
		ConfigFilename: "/etc/rfoutlet/config.yml",
	}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the frontend for controlling outlets",
		Long:  "The serve command starts a server which serves the frontend and connects clients through websockets for controlling outlets via web interface.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return options.Run()
		},
	}

	options.AddFlags(cmd)

	return cmd
}

type ServeOptions struct {
	config.Config
	ConfigFilename string
}

func (o *ServeOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFilename, "config", o.ConfigFilename, "config filename")
	cmd.Flags().StringVar(&o.Config.StateFile, "state-file", o.Config.StateFile, "state filename")
	cmd.Flags().StringVar(&o.Config.ListenAddress, "listen-address", o.Config.ListenAddress, "listen address")
	cmd.Flags().UintVar(&o.Config.TransmitPin, "transmit-pin", o.Config.TransmitPin, "gpio pin to transmit on")
}

func (o *ServeOptions) Run() error {
	config, err := config.LoadWithDefaults(o.ConfigFilename)
	if err != nil {
		return err
	}

	err = mergo.Merge(config, o.Config, mergo.WithOverride)
	if err != nil {
		return err
	}

	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}
	defer chip.Close()

	transmitter, err := gpio.NewTransmitter(chip, int(config.TransmitPin))
	if err != nil {
		return err
	}
	defer transmitter.Close()

	registry := outlet.NewRegistry()

	err = registry.RegisterGroups(config.BuildOutletGroups()...)
	if err != nil {
		return err
	}

	if config.StateFile != "" {
		outletState, err := state.Load(config.StateFile)
		if err == nil {
			outletState.Apply(registry.GetOutlets())
		} else if !os.IsNotExist(err) {
			return err
		}

		defer func() {
			outletState := state.Collect(registry.GetOutlets())

			err := state.Save(config.StateFile, outletState)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	stopCh := make(chan struct{})
	commandQueue := make(chan command.Command)

	hub := websocket.NewHub()

	controller := controller.Controller{
		Registry:     registry,
		Switcher:     outlet.NewSwitch(transmitter),
		Broadcaster:  hub,
		CommandQueue: commandQueue,
	}

	sched := scheduler.New(registry, commandQueue)

	go controller.Run(stopCh)
	go sched.Run(stopCh)
	go hub.Run(stopCh)
	go handleSignals(stopCh)

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", handler.Redirect("/app"))
	router.GET("/healthz", handler.Healthz)
	router.GET("/ws", handler.Websocket(hub, commandQueue))
	router.StaticFS("/app", packr.NewBox(webDir))

	return listenAndServe(stopCh, router, config.ListenAddress)
}

func listenAndServe(stopCh <-chan struct{}, handler http.Handler, addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	<-stopCh

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func handleSignals(stopCh chan struct{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("received signal, terminating...")
	close(stopCh)
}
