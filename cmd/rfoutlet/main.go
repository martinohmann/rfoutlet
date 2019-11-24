// The rfoutlet command starts a server which serves the frontend and connects
// clients through websockets for controlling outlets via web interface.
//
// Available command line flags:
//
//  -config string
//        config filename (default "/etc/rfoutlet/config.yml")
//  -gpio-pin uint
//        gpio pin to transmit on (default -1)
//  -listen-address string
//        listen address
//  -state-file string
//        state filename
package main

import (
	"context"
	"flag"
	"fmt"
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
)

const (
	webDir                = "../../web/build"
	defaultConfigFilename = "/etc/rfoutlet/config.yml"
)

var (
	configFilename = flag.String("config", defaultConfigFilename, "config filename")
	stateFilename  = flag.String("state-file", "", "state filename")
	listenAddress  = flag.String("listen-address", "", "listen address")
	gpioPin        = flag.Int("gpio-pin", -1, "gpio pin to transmit on")
	usage          = func() {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func init() {
	flag.Usage = usage
}

func exitError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	config, err := config.Load(*configFilename)
	exitError(err)

	if *gpioPin >= 0 {
		config.GpioPin = uint(*gpioPin)
	}

	if *listenAddress != "" {
		config.ListenAddress = *listenAddress
	}

	if *stateFilename != "" {
		config.StateFile = *stateFilename
	}

	manager := outlet.NewManager(state.NewHandler(config.StateFile))
	defer manager.SaveState()

	err = outlet.RegisterFromConfig(manager, config)
	exitError(err)

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

	listenAndServe(router, config.ListenAddress)
}

func listenAndServe(handler http.Handler, addr string) {
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
