package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/martinohmann/rfoutlet/internal/api"
	"github.com/martinohmann/rfoutlet/internal/gpio"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

const (
	webDir                = "../../app/build"
	defaultListenAddress  = "0.0.0.0:3333"
	defaultConfigFilename = "config.yml"
)

var (
	configFilename = flag.String("config", defaultConfigFilename, "config filename")
	listenAddress  = flag.String("listen-address", defaultListenAddress, "listen address")
	gpioPin        = flag.Int("gpio-pin", gpio.DefaultGpioPin, "gpio pin to transmit on")
)

func main() {
	flag.Parse()

	outletConfig := outlet.ReadConfig(*configFilename)
	transmitter, err := gpio.NewTransmitter(*gpioPin)
	if err != nil {
		panic(err)
	}

	defer transmitter.Close()

	box := packr.NewBox(webDir)

	api := api.New(outletConfig, transmitter)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags|log.Lshortfile)

	router := http.NewServeMux()

	router.Handle("/", http.FileServer(box))
	router.HandleFunc("/api/status", api.HandleStatusRequest)
	router.HandleFunc("/api/outlet_group/", api.ValidateRequest(api.HandleOutletGroupRequest))
	router.HandleFunc("/api/outlet/", api.ValidateRequest(api.HandleOutletRequest))

	server := &http.Server{
		Addr:    *listenAddress,
		Handler: logging(logger)(router),
	}

	logger.Printf("Listening on %s\n", *listenAddress)

	server.ListenAndServe()
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
