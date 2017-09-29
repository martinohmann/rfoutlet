package main

import (
	"log"
	"net/http"
	"os"
	"rf-outlet/internal"
)

const (
	webDir = "frontend/build"
)

func main() {
	configFilename := os.Getenv("RF_CONFIG")
	if configFilename == "" {
		configFilename = internal.DefaultConfigFilename
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "0.0.0.0:3000"
	}

	config := internal.ReadConfig(configFilename)
	config.Print()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	api := internal.NewAPI(config)
	http.HandleFunc("/api/outlet_group/", api.ValidateRequest(api.HandleOutletGroupRequest))
	http.HandleFunc("/api/outlet/", api.ValidateRequest(api.HandleOutletRequest))

	log.Printf("Listening on %s...\n", listenAddress)
	http.ListenAndServe(listenAddress, nil)
}
