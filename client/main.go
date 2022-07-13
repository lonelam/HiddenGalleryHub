package main

import (
	"HiddenGalleryHub/client/processor"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "localhost:5555", "http service address")
var rootDir = flag.String("root", ".", "the root dir to search images")

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/"}
	log.Printf("connecting to %s", u.String())

	c := processor.CreateWsConnection(u.String(), *rootDir)

	done := c.StartUp()

	select {
	case <-done:
	case <-interrupt:
	}
}
