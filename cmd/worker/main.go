package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/OnYyon/gRPCCalculator/internal/transport/grpc/worker"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var workers sync.Map
	for i := 0; i < 3; i++ {
		wg.Add(1)
		fmt.Println("worker start - ", i)
		go func(i int) {
			defer wg.Done()
			worker, err := worker.NewWorker()
			if err != nil {
				fmt.Println(err)
				return
			}
			workers.Store(i, worker)

			if err := worker.Run(); err != nil {
				fmt.Println("error")
			}
		}(i)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancel()

	workers.Range(func(key, value interface{}) bool {
		if w, ok := value.(*worker.Worker); ok {
			w.Stop()
		}
		return true
	})

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	log.Println("Workers stopped gracefully")
}
