package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/buzzsurfr/exorcism/http"

	log "k8s.io/klog/v2"
)

var (
	tlscert, tlskey, port string
	version               bool
)

func main() {
	flag.StringVar(&tlscert, "tlscert", "/etc/certs/tls.crt", "Path to the TLS certificate")
	flag.StringVar(&tlskey, "tlskey", "/etc/certs/tls.key", "Path to the TLS key")
	flag.StringVar(&port, "port", "8443", "The port to listen")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Println("Version built on 2022-10-26")
		os.Exit(1)
	}

	server := http.NewServer(port)
	go func() {
		if err := server.ListenAndServeTLS(tlscert, tlskey); err != nil {
			log.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	log.Infof("Server running in port: %s", port)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Infof("Shutdown gracefully...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}
