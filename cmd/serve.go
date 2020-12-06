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
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/statedrift"
	"github.com/martinohmann/rfoutlet/internal/timeswitch"
	"github.com/martinohmann/rfoutlet/internal/websocket"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
			return options.Run(cmd)
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
	cmd.Flags().StringVar(&o.ConfigFilename, "config", o.ConfigFilename, "path to the outlet config file")
	cmd.Flags().StringVar(&o.StateFile, "state-file", o.StateFile, "path to the file where outlet state and schedule should be stored")
	cmd.Flags().StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "address to serve the web app on")
	cmd.Flags().BoolVar(&o.DetectStateDrift, "detect-state-drift", o.DetectStateDrift, "detect state drift (e.g. if an outlet was switched via the phyical remote instead of rfoutlet)")
	cmd.Flags().UintVar(&o.GPIO.TransmitPin, "transmit-pin", o.GPIO.TransmitPin, "gpio pin to transmit rf codes on")
	cmd.Flags().UintVar(&o.GPIO.ReceivePin, "receive-pin", o.GPIO.ReceivePin, "gpio pin to receive rf codes on (this is used by the state drift detector)")
	cmd.Flags().IntVar(&o.GPIO.TransmissionCount, "transmission-count", o.GPIO.TransmissionCount, "number of times a code should be transmitted in a row. The higher the value, the more likely it is that an outlet actually received the code")
}

func (o *ServeOptions) Run(cmd *cobra.Command) error {
	cfg, err := config.LoadWithDefaults(o.ConfigFilename)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	err = mergo.Merge(cfg, o.Config, mergo.WithOverride)
	if err != nil {
		return fmt.Errorf("failed to merge config values: %v", err)
	}

	log.Debugf("merged config values: %#v", cfg)

	registry := outlet.NewRegistry()

	err = registry.RegisterGroups(cfg.BuildOutletGroups()...)
	if err != nil {
		return fmt.Errorf("failed to register outlet groups: %v", err)
	}

	device, err := openGPIODevice(cmd)
	if err != nil {
		return err
	}
	defer device.Close()

	transmitter, err := gpio.NewTransmitter(device.Chip, int(cfg.GPIO.TransmitPin), gpio.TransmissionCount(cfg.GPIO.TransmissionCount))
	if err != nil {
		return fmt.Errorf("failed to create gpio transmitter: %v", err)
	}
	defer transmitter.Close()

	if cfg.StateFile != "" {
		log := log.WithField("stateFile", cfg.StateFile)

		stateFile := outlet.NewStateFile(cfg.StateFile)

		log.Debug("loading outlet states")

		err := stateFile.ReadBack(registry.GetOutlets())
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to load outlet states: %v", err)
		}

		defer func() {
			log.Info("saving outlet states")

			err := stateFile.WriteOut(registry.GetOutlets())
			if err != nil {
				log.Errorf("failed to save state: %v", err)
			}
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopCh := ctx.Done()
	commandQueue := make(chan command.Command)

	if cfg.DetectStateDrift {
		receiver, err := gpio.NewReceiver(device.Chip, int(cfg.GPIO.ReceivePin))
		if err != nil {
			return fmt.Errorf("failed to create gpio receiver: %v", err)
		}
		defer receiver.Close()

		detector := statedrift.NewDetector(registry, receiver, commandQueue)

		go detector.Run(stopCh)
	}

	hub := websocket.NewHub()

	controller := controller.Controller{
		Registry:     registry,
		Switcher:     outlet.NewSwitch(transmitter),
		Broadcaster:  hub,
		CommandQueue: commandQueue,
	}

	timeSwitch := timeswitch.New(registry, commandQueue)

	go handleSignals(cancel)
	go controller.Run(stopCh)
	go timeSwitch.Run(stopCh)
	go hub.Run(stopCh)

	router := setupRouter(hub, commandQueue)

	return listenAndServe(stopCh, router, cfg.ListenAddress)
}

func setupRouter(hub *websocket.Hub, commandQueue chan<- command.Command) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), cors.Default())
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/app") })
	r.GET("/ws", websocket.Handler(hub, commandQueue))
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.StaticFS("/app", packr.NewBox(webDir))

	return r
}

func listenAndServe(stopCh <-chan struct{}, handler http.Handler, addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Infof("listening on %s", addr)

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

func handleSignals(cancel func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	log.WithField("signal", sig).Info("received signal, shutting down...")
	cancel()
}
