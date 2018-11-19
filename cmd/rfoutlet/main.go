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
	ctx "context"
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
	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/martinohmann/rfoutlet/internal/scheduler"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

const (
	webDir                = "../../app/build"
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

func main() {
	flag.Parse()

	config, err := config.Load(*configFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *gpioPin >= 0 {
		config.GpioPin = uint(*gpioPin)
	}

	transmitter, err := gpio.NewTransmitter(config.GpioPin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer transmitter.Close()

	if *listenAddress != "" {
		config.ListenAddress = *listenAddress
	}

	if *stateFilename != "" {
		config.StateFile = *stateFilename
	}

	s, err := state.Load(config.StateFile)
	if err != nil {
		s = state.New()
	}

	ctx, err := context.New(config, s)

	control := control.New(ctx, transmitter)
	scheduler := scheduler.New(ctx, control, 10*time.Second)

	scheduler.Start()

	defer scheduler.Stop()

	api := api.New(ctx, control)

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", handler.Redirect("/app"))
	router.GET("/healthz", handler.Healthz)
	router.StaticFS("/app", packr.NewBox(webDir))

	apiRoutes := router.Group("/api")
	apiRoutes.GET("/status", api.StatusRequestHandler)
	apiRoutes.POST("/outlet", api.OutletRequestHandler)
	apiRoutes.POST("/outlet_group", api.GroupRequestHandler)
	apiRoutes.PUT("/outlet/schedule", api.IntervalRequestHandler)
	apiRoutes.POST("/outlet/schedule", api.IntervalRequestHandler)
	apiRoutes.DELETE("/outlet/schedule", api.IntervalRequestHandler)

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

	ctx, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
