package worker

import (
	"context"
	"fmt"
	"log"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Worker struct {
	conn   *grpc.ClientConn
	client proto.OrchestratorClient
	stream proto.Orchestrator_TaskStreamClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorker() (*Worker, error) {
	ctx, cancel := context.WithCancel(context.Background())

	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := proto.NewOrchestratorClient(conn)
	stream, err := client.TaskStream(ctx)
	if err != nil {
		conn.Close()
		cancel()
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return &Worker{
		conn:   conn,
		client: client,
		stream: stream,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (w *Worker) Run() error {
	defer w.cleanup()

	for {
		select {
		case <-w.ctx.Done():
			return nil
		default:
			task, err := w.stream.Recv()
			if err != nil {
				return fmt.Errorf("receive error: %w", err)
			}

			result := processTask(task)

			if err := w.stream.Send(result); err != nil {
				return fmt.Errorf("send error: %w", err)
			}
		}
	}
}

func (w *Worker) Stop() {
	w.cancel()
}

func (w *Worker) cleanup() {
	if w.stream != nil {
		if err := w.stream.CloseSend(); err != nil {
			log.Printf("Error closing stream: %v", err)
		}
	}
	if w.conn != nil {
		if err := w.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}
}

func processTask(task *proto.Task) *proto.Task {
	task.ID = "modifed"
	return task
}
