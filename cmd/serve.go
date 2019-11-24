package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/scheduler"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/spf13/cobra"
)

const webDir = "../web/build"

func NewServeCommand() *cobra.Command {
	options := &ServeOptions{
		ConfigFilename: "/etc/rfoutlet/config.yml",
		GpioPin:        -1,
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
	ConfigFilename string
	StateFilename  string
	ListenAddress  string
	GpioPin        int
}

func (o *ServeOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFilename, "config", o.ConfigFilename, "config filename")
	cmd.Flags().StringVar(&o.StateFilename, "state-file", o.StateFilename, "state filename")
	cmd.Flags().StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "listen address")
	cmd.Flags().IntVar(&o.GpioPin, "gpio-pin", o.GpioPin, "gpio pin to transmit on")
}

func (o *ServeOptions) Run() error {
	config, err := config.Load(o.ConfigFilename)
	if err != nil {
		return err
	}

	if o.GpioPin >= 0 {
		config.GpioPin = uint(o.GpioPin)
	}

	if o.ListenAddress != "" {
		config.ListenAddress = o.ListenAddress
	}

	if o.StateFilename != "" {
		config.StateFile = o.StateFilename
	}

	manager := outlet.NewManager(state.NewHandler(config.StateFile))
	defer manager.SaveState()

	err = outlet.RegisterFromConfig(manager, config)
	if err != nil {
		return err
	}

	manager.LoadState()

	transmitter := gpio.NewTransmitter(config.GpioPin)
	defer transmitter.Close()

	switcher := outlet.NewSwitch(transmitter)
	hub := control.NewHub()
	control := control.New(manager, switcher, hub)
	scheduler := scheduler.New(control)

	for _, o := range manager.Outlets() {
		scheduler.Register(o)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", handler.Redirect("/app"))
	router.GET("/healthz", handler.Healthz)
	router.GET("/ws", handler.Websocket(hub, control))
	router.StaticFS("/app", packr.NewBox(webDir))

	return listenAndServe(router, config.ListenAddress)
}

func listenAndServe(handler http.Handler, addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
