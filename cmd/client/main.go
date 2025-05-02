package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerClient struct {
	client proto.OrchestratorClient
	stream proto.Orchestrator_TaskStreamClient
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorkerClient(conn *grpc.ClientConn) *WorkerClient {
	ctx, cancel := context.WithCancel(context.Background())
	client := proto.NewOrchestratorClient(conn)
	stream, err := client.TaskStream(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	return &WorkerClient{
		client: client,
		stream: stream,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (wc *WorkerClient) Start() {
	// Горутина для чтения ответов от сервера
	wc.wg.Add(1)
	go func() {
		defer wc.wg.Done()
		for {
			resp, err := wc.stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				return
			}
			log.Printf("Received response for task %v", resp)
		}
	}()

	// Горутина для отправки задач
	wc.wg.Add(1)
	go func() {
		defer wc.wg.Done()
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				taskID := fmt.Sprintf("task-%d", rand.Intn(1000))
				task := &proto.Task{
					ID: taskID,
				}

				if err := wc.stream.Send(task); err != nil {
					log.Printf("Send error: %v", err)
					return
				}
				log.Printf("Sent task: %s", taskID)

			case <-wc.ctx.Done():
				return
			}
		}
	}()
}

func (wc *WorkerClient) Stop() {
	wc.cancel()
	wc.wg.Wait()
	if err := wc.stream.CloseSend(); err != nil {
		log.Printf("Failed to close stream: %v", err)
	}
}

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	worker := NewWorkerClient(conn)
	worker.Start()

	// Ожидание сигнала для завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	worker.Stop()
	log.Println("Worker stopped gracefully")
}
