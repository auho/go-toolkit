package main

import (
	"log"

	"dns/mdns/publisher"
)

func main() {
	err := publisher.RunFile()
	if err != nil {
		log.Fatal(err)
	}
}
