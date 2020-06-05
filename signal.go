package vaali

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Signal struct {
	once sync.Once
}

func (s *Signal) GetContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	s.once.Do(func() {
		sigChannel := make(chan os.Signal, 2)
		shutdownSignals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

		signal.Notify(sigChannel, shutdownSignals...)

		go func() {
			<-sigChannel
			cancel()
			<-sigChannel
			os.Exit(1) // second signal. Exit directly.
		}()
	})

	return ctx
}
