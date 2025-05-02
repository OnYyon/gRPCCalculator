package orchestrator

import (
	"fmt"
	"io"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

type serverAPI struct {
	proto.UnimplementedOrchestratorServer
	// tasks chan *proto.Task
	// results map[string][]float64
	// mu sync.Mutex
}

func RegisterOrchestratorServer(gRPC *grpc.Server) {
	proto.RegisterOrchestratorServer(gRPC, &serverAPI{})
}

func (s *serverAPI) TaskStream(stream grpc.BidiStreamingServer[proto.Task, proto.Task]) error {
	for {
		task, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// Обработка задачи
		go func(t *proto.Task) {
			// ... обработка задачи ...
			fmt.Println(t)
			// Отправка ответа
			resp := &proto.Task{
				ID: t.ID,
			}
			stream.Send(resp)
		}(task)
	}
}

// func RegisterOrchestratorServer(gRPC *grpc.Server) *serverAPI {
// 	s := &serverAPI{
// 		tasks: make(chan *proto.Task, 3),
// 	}
// 	proto.RegisterOrchestratorServer(gRPC, s)
// 	return s
// }

// func (s *serverAPI) GetTask(ctx context.Context, _ *proto.TypeNil) (*proto.Task, error) {
// 	select {
// 	case task := <-s.tasks:
// 		return task, nil
// 	default:
// 		return nil, errors.New("no tasks")
// 	}
// }

// // NOTE: for tests may be
// func (s *serverAPI) AddTask(task *proto.Task) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	go func() { s.tasks <- task }()
// }
