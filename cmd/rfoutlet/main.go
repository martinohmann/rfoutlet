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
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/martinohmann/rfoutlet/internal/api"
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

	outletConfig, err := outlet.ReadConfig(*configFilename)
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
		Handler: cors("*")(logging(logger)(router)),
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

func cors(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	}
}
