package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aren55555/future-hw/api"
	"github.com/aren55555/future-hw/data/mem"
)

var (
	port         = flag.Int("port", 8080, "the port to listen on")
	fileLocation = flag.String("location", "", "the location of the seed data JSON file")
)

func main() {
	flag.Parse()

	memDataStore := mem.New()
	memDataStore.Seed(*fileLocation)
	apiHandler := api.New(memDataStore)

	http.Handle("/api", apiHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

	// TODO: could dump the file once the server has halted via a graceful shutdown
}
