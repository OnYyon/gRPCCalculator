package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Worker struct {
	client proto.OrchestratorClient
	stream proto.Orchestrator_TaskStreamClient
}

func NewWorkerClient(conn *grpc.ClientConn) *Worker {
	client := proto.NewOrchestratorClient(conn)
	stream, err := client.TaskStream(context.Background())
	if err != nil {
		log.Fatalf("Failed to establish stream: %v", err)
	}
	return &Worker{
		client: client,
		stream: stream,
	}
}

func (w *Worker) Start() {
	// Обработка задач
	go func() {
		for {
			task, err := w.stream.Recv()
			if err != nil {
				log.Printf("Failed to receive task: %v", err)
				return
			}
			task.ID = "modifed"
			// Отправка результата
			w.stream.Send(task)
		}
	}()
}

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	worker := NewWorkerClient(conn)
	for i := 0; i < 3; i++ {
		fmt.Println("worker start - ", i)
		go worker.Start()
	}

	// Ожидание сигнала для завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Worker stopped gracefully")
}
