package vaali

import (
	"context"
	"log"
	"runtime"
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
	MemLowerBound int64
	MemUpperBound int64
}

func (s *Synchronizer) Start(stop <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			s.wg.Wait()
			for _, cancel := range s.cancellations {
				cancel()
			}

			return
		case <-ticker.C:
			log.Printf("Number of active goroutines: %d\n", runtime.NumGoroutine())
			s.adjust(stop)
		}
	}
}

func (s *Synchronizer) adjust(stop <-chan struct{}) {
	log.Printf("Memory usage: %d\n", memUsage())

	for memUsage() < uint64(s.MemLowerBound) {
		ctx, cancel := context.WithCancel(context.Background())
		s.cancellations = append(s.cancellations, cancel)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			s.Runner.Run(ctx, stop)
		}()

		time.Sleep(100 * time.Millisecond)
	}

	for memUsage() > uint64(s.MemUpperBound) {
		if len(s.cancellations) == 0 {
			return
		}
		cancel := s.cancellations[len(s.cancellations)-1]
		s.cancellations = s.cancellations[:len(s.cancellations)-1]

		cancel()

		time.Sleep(100 * time.Millisecond)
	}
}

func (s *Synchronizer) cancelGoroutines() {
}

func memUsage() uint64 {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	return byteToMegabyte(memStats.Sys)
}

func byteToMegabyte(b uint64) uint64 {
	return b / 1024 / 1024
}
