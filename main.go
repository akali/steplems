package main

import (
	"log"
)

func main() {
	app, err := NewWireApplication()

	if err != nil {
		log.Fatalf("Failed to create wire app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start wire app: %v", err)
	}
}
