package main

import (
	"log"
	"os"

	"github.com/narahari92/vaali"
	"github.com/narahari92/vaali/rand"
)

func main() {
	run()
}

func run() {
	signal := vaali.Signal{}
	ctx := signal.GetContext()
	synchronizer := vaali.Synchronizer{
		Runner: rand.Generator(),
	}

	if err := synchronizer.Start(ctx.Done()); err != nil {
		log.Printf("program failed, err: %v", err)
		os.Exit(1)
	}

	log.Println("program completed")
}
