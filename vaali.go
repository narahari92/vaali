package vaali

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Runnable interface {
	Run(context.Context, <-chan struct{})
}

type RunnableFunc func(context.Context, <-chan struct{})

func (rf RunnableFunc) Run(ctx context.Context, stop <-chan struct{}) {
	rf(ctx, stop)
}

type Synchronizer struct {
	wg            sync.WaitGroup
	cancellations []context.CancelFunc
	Runner        Runnable
}

func (s *Synchronizer) Start(stop <-chan struct{}) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			s.wg.Wait()
			s.cancelGoroutines()

			return nil
		case <-ticker.C:
			log.Printf("Number of active goroutines: %d\n", runtime.NumGoroutine())
			err := s.adjust(stop)
			if err != nil {
				s.cancelGoroutines()
				return err
			}
		}
	}
}

func (s *Synchronizer) adjust(stop <-chan struct{}) error {
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

	log.Printf("Memory usage: %d\n", memUsage())
	for memUsage() < uint64(memLowerBound) {
		ctx, cancel := context.WithCancel(context.Background())
		s.cancellations = append(s.cancellations, cancel)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			s.Runner.Run(ctx, stop)
		}()

		time.Sleep(100 * time.Millisecond)
	}

	for memUsage() > uint64(memUpperBound) {
		if len(s.cancellations) == 0 {
			return nil
		}
		cancel := s.cancellations[len(s.cancellations)-1]
		s.cancellations = s.cancellations[:len(s.cancellations)-1]

		cancel()

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (s *Synchronizer) cancelGoroutines() {
	for _, cancel := range s.cancellations {
		cancel()
	}
}

func memUsage() uint64 {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	return byteToMegabyte(memStats.Sys)
}

func byteToMegabyte(b uint64) uint64 {
	return b / 1024 / 1024
}
