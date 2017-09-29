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
	http.HandleFunc("/api/", serveApi)

	log.Printf("Listening on %s...\n", listenAddress)
	http.ListenAndServe(listenAddress, nil)
}

func serveApi(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("api call"))

	log.Println(r.RequestURI)
}
