package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"proxy/internal"
	"strconv"
	"syscall"
)

func main() {
	// Read config
	if len(os.Args) < 2 || os.Args[1] == "" {
		panic("Config file path not provided")
	}
	err := internal.ReadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Initialize logger
	internal.InitLogger()

	// Start server
	exitSig := make(chan os.Signal, 1)
	signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(internal.Config.Http_port), http.HandlerFunc(internal.ProxyHandler))
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("%v", err)
			internal.CloseLogger()
			os.Exit(1)
		}
	}()

	log.Printf("Listening HTTP on port %d...\n", internal.Config.Http_port)
	<-exitSig
	log.Printf("Shutting down...")
	internal.CloseLogger()
}
