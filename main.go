package main

import "log"

func main() {
	app, err := NewWireApplication()

	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}
