package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/narahari92/vaali"
	"github.com/narahari92/vaali/rand"
)

func main() {
	if err := run(); err != nil {
		log.Printf("program failed, err: %v", err)
		os.Exit(1)
	}

	log.Println("program completed")
}

func run() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("program expects 2 integer arguments for lower bound memory and upper bound memory")
	}

	memLowerBound, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert lower bound memory argument %s into int64, err: %v", os.Args[0], err)
	}

	memUpperBound, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert upper bound memory argument %s into int64, err: %v", os.Args[1], err)
	}

	signal := vaali.Signal{}
	ctx := signal.GetContext()
	synchronizer := vaali.Synchronizer{
		Runner:        rand.Generator(),
		MemLowerBound: memLowerBound,
		MemUpperBound: memUpperBound,
	}

	synchronizer.Start(ctx.Done())

	return nil
}
