package cmd

import (
	"context"
	"fmt"
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
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/internal/timeswitch"
	"github.com/martinohmann/rfoutlet/internal/websocket"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	log "github.com/sirupsen/logrus"
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
		RunE: func(_ *cobra.Command, _ []string) error {
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
	cmd.Flags().StringVar(&o.StateFile, "state-file", o.StateFile, "state filename")
	cmd.Flags().StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "listen address")
	cmd.Flags().UintVar(&o.TransmitPin, "transmit-pin", o.TransmitPin, "gpio pin to transmit on")
}

func (o *ServeOptions) Run() error {
	cfg, err := config.LoadWithDefaults(o.ConfigFilename)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	err = mergo.Merge(cfg, o.Config, mergo.WithOverride)
	if err != nil {
		return fmt.Errorf("failed to merge config values: %v", err)
	}

	log.WithField("config", fmt.Sprintf("%#v", cfg)).Debug("merged config values")

	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return fmt.Errorf("failed to open gpio device: %v", err)
	}
	defer chip.Close()

	transmitter, err := gpio.NewTransmitter(chip, int(cfg.TransmitPin))
	if err != nil {
		return fmt.Errorf("failed to create gpio transmitter: %v", err)
	}
	defer transmitter.Close()

	registry := outlet.NewRegistry()

	err = registry.RegisterGroups(cfg.BuildOutletGroups()...)
	if err != nil {
		return fmt.Errorf("failed to registry outlet groups: %v", err)
	}

	if cfg.StateFile != "" {
		log := log.WithField("stateFile", cfg.StateFile)

		log.Debug("loading outlet states")

		outletState, err := state.Load(cfg.StateFile)
		if err == nil {
			outletState.Apply(registry.GetOutlets())
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to load outlet states: %v", err)
		}

		defer func() {
			log.Info("saving outlet states")

			outletState := state.Collect(registry.GetOutlets())

			err := state.Save(cfg.StateFile, outletState)
			if err != nil {
				log.Errorf("failed to save state: %v", err)
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

	timeSwitch := timeswitch.New(registry, commandQueue)

	go handleSignals(stopCh)
	go controller.Run(stopCh)
	go timeSwitch.Run(stopCh)
	go hub.Run(stopCh)

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", handler.Redirect("/app"))
	r.GET("/ws", handler.Websocket(hub, commandQueue))
	r.GET("/healthz", handler.Healthz)
	r.StaticFS("/app", packr.NewBox(webDir))

	return listenAndServe(stopCh, r, cfg.ListenAddress)
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

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func handleSignals(stopCh chan struct{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("received signal, terminating...")
	close(stopCh)
}
