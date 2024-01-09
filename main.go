package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var build = "develop"

func main() {
	log.Println("starting service", build)
	defer log.Println("service ended")

	//Controlled shutdown, SIGINT is for CTRL+C
	// SIGTERM is what Kubernetes is using to terminate the app
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stopping service")
}
