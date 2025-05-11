package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	services "github.com/OnYyon/gRPCCalculator/internal/services/calculate"
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

	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 3; i++ {
		conn, err = grpc.NewClient(
			"localhost:8080",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			break
		}

		waitTime := time.Second * time.Duration(i+1)
		log.Printf("Connection attempt %d failed, retrying in %v: %v", i+1, waitTime, err)
		time.Sleep(waitTime)
	}

	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect after 3 attempts: %w", err)
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
			fmt.Printf("get task: %v %v %v\n", task.Arg1, task.Operator, task.Arg2)

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(task.Timeout))
			defer cancel()

			result, err := services.ProcessTask(ctx, task)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Println("timeout")
					result.DescErr = "timeout"
					result.Completed = false
					result.Err = true
				} else {
					result.Completed = false
					result.Err = true
					result.DescErr = fmt.Sprint(err)
				}
			}

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
