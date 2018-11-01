package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/martinohmann/rfoutlet/internal/api"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

const (
	defaultWebDir         = "app/build"
	defaultListenAddress  = "0.0.0.0:3333"
	defaultConfigFilename = "config.yml"
)

var (
	webDir         = flag.String("web-dir", defaultWebDir, "web directory")
	configFilename = flag.String("config", defaultConfigFilename, "config filename")
	listenAddress  = flag.String("listen-address", defaultListenAddress, "listen address")
)

func main() {
	flag.Parse()

	config := outlet.ReadConfig(*configFilename)

	api := api.New(config)

	http.Handle("/", http.FileServer(http.Dir(*webDir)))
	http.HandleFunc("/api/status", api.HandleStatusRequest)
	http.HandleFunc("/api/outlet_group/", api.ValidateRequest(api.HandleOutletGroupRequest))
	http.HandleFunc("/api/outlet/", api.ValidateRequest(api.HandleOutletRequest))

	log.Printf("Listening on %s\n", *listenAddress)

	config.Print()

	http.ListenAndServe(*listenAddress, nil)
}
