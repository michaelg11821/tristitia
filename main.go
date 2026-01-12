package main

import (
	"log"
)

func main() {
	engine, err := createEngine()
	if err != nil {
		log.Fatal(err)
	}

	for _, ncode := range engine.NCodes {
		go engine.watch(ncode)
	}

	select {}
}
