package main

import (
	"log"

	"example.com/hello-goui/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
