// The rfoutlet command starts a server which serves the frontend and api for
// controlling outlets via web interface.
//
// Available command line flags:
//
//  -config string
//        config filename (default "/etc/rfoutlet/config.yml")
//  -gpio-pin uint
//        gpio pin to transmit on (default 17)
//  -listen-address string
//        listen address (default "0.0.0.0:3333")
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
	"github.com/martinohmann/rfoutlet/internal/api"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

const (
	webDir                = "../../app/build"
	defaultListenAddress  = "0.0.0.0:3333"
	defaultConfigFilename = "/etc/rfoutlet/config.yml"
)

var (
	configFilename = flag.String("config", defaultConfigFilename, "config filename")
	stateFilename  = flag.String("state-file", "", "state filename")
	listenAddress  = flag.String("listen-address", defaultListenAddress, "listen address")
	gpioPin        = flag.Uint("gpio-pin", gpio.DefaultTransmitPin, "gpio pin to transmit on")
	usage          = func() {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func init() {
	flag.Usage = usage
}

func main() {
	flag.Parse()

	config, err := outlet.ReadConfig(*configFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	transmitter, err := gpio.NewTransmitter(*gpioPin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer transmitter.Close()

	var stateManager outlet.StateManager

	if *stateFilename != "" {
		stateFile, err := os.OpenFile(*stateFilename, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		stateManager = outlet.NewStateManager(stateFile)
	} else {
		stateManager = outlet.NewNullStateManager()
	}

	defer stateManager.Close()

	control := outlet.NewControl(config, stateManager, transmitter)

	if err := control.RestoreState(); err != nil {
		log.Printf("error while restoring state: %s\n", err)
	}

	api := api.New(control)

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", handler.Redirect("/app"))
	router.GET("/healthz", handler.Healthz)
	router.StaticFS("/app", packr.NewBox(webDir))

	apiRoutes := router.Group("/api")
	apiRoutes.GET("/status", api.StatusRequestHandler)
	apiRoutes.POST("/outlet", api.OutletRequestHandler)
	apiRoutes.POST("/outlet_group", api.OutletGroupRequestHandler)

	listenAndServe(router, *listenAddress)
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
