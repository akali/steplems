package main

import (
	"flag"
	"log"
)

var command = flag.String("command", "runbot", "command to execute")

func main() {
	flag.Parse()

	app, err := NewWireApplication()

	if err != nil {
		log.Fatalf("Failed to create wire app: %v", err)
	}

	if err := app.Start(*command); err != nil {
		log.Fatalf("Failed to start wire app: %v", err)
	}
}
