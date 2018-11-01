package main

import (
	"flag"
	"log"
	"net/http"
	"os"

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

	logger := log.New(os.Stdout, "http: ", log.LstdFlags|log.Lshortfile)

	router := http.NewServeMux()

	router.Handle("/", http.FileServer(http.Dir(*webDir)))
	router.HandleFunc("/api/status", api.HandleStatusRequest)
	router.HandleFunc("/api/outlet_group/", api.ValidateRequest(api.HandleOutletGroupRequest))
	router.HandleFunc("/api/outlet/", api.ValidateRequest(api.HandleOutletRequest))

	server := &http.Server{
		Addr:    *listenAddress,
		Handler: logging(logger)(router),
	}

	logger.Printf("Listening on %s\n", *listenAddress)

	config.Print()

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
