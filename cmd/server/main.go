package main

import (
	"log"
	"server/internal/transport/rest"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := rest.RunRest()
	if err != nil {
		return err
	}
	return nil
}
