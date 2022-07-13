package main

import (
	"HiddenGalleryHub/client/processor"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "laizn.com", "http service address")
var rootDir = flag.String("root", ".", "the root dir to search images")
var name = flag.String("name", "home", "the name to identify machines")

func main() {
	flag.Parse()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/ws/"}
	log.Printf("connecting to %s", u.String())

	c := processor.CreateWsConnection(u.String(), *rootDir, *name)

	done := c.StartUp()

	select {
	case <-done:
	case <-interrupt:
	}
}
