package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

var pingInterval = flag.Duration("interval", 30*time.Second, "Ping interval")
var pingURL = flag.String("url", "", "Destination Ping HTTP URL")

func main() {
	log.Println("Starting off the application")
	flag.Parse()
	if *pingURL == "" {
		log.Fatalln("invalid input param, please specify a Ping URL")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ping, err := NewPinger(*pingInterval, *pingURL)
	if err != nil {
		log.Fatalf("failed to init pinger: %s", err)
	}
	ping.cancel = cancel

	go ping.Start(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	ping.Stop()
	log.Println("Stopping on the application")
}
