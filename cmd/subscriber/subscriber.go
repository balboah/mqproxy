package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/apcera/nats"
)

type request struct {
	Method string
	Path   string
	Body   string
}

type response struct {
	Code int
	Body string
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	name := flag.String("name", "service", "The service to publish to")
	url := flag.String("nats", "nats://localhost:4222", "The nats URL")
	flag.Parse()

	nc, err := nats.Connect(*url)
	if err != nil {
		log.Fatal(err)
	}
	c, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Requests
	log.Println("Reading requests")
	c.Subscribe(*name, func(subj, reply string, req *request) {
		c.Publish(reply, &response{200, "Hello world!"})
	})

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
