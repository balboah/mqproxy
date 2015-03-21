package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

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

// natsProxy proxies requests via a NATS server.
type natsProxy struct {
	conn *nats.EncodedConn

	// Subject is where to send requests and receive replies from.
	subj string
}

// NatsHandler implements the http.Handler via a Nats connection.
func NatsHandler(conn *nats.EncodedConn, subj string) http.Handler {
	return &natsProxy{conn: conn, subj: subj}
}

func (p *natsProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var httpRes response
	if err := p.conn.Request(
		p.subj, &request{req.Method, req.RequestURI, ""}, &httpRes, 30*time.Second,
	); err != nil {
		log.Println(err)
		rw.WriteHeader(500)
		fmt.Fprintln(rw, err.Error())
		return
	}

	rw.WriteHeader(httpRes.Code)
	fmt.Fprint(rw, httpRes.Body)
}

func normal(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "Hello world!")
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

	log.Println("Listening...")
	http.Handle("/nats", NatsHandler(c, *name))
	http.Handle("/normal", http.HandlerFunc(normal))
	http.ListenAndServe(":8080", nil)
}
