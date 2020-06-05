package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"gitlab.eng.vmware.com/hnarahari/vaali"
	"gitlab.eng.vmware.com/hnarahari/vaali/rand"
)

func main() {
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("profiling server failed, err: %v", err)
			os.Exit(1)
		}

		log.Println("profiling server closed")
	}()

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
