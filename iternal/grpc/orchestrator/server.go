package orchestrator

import (
	"fmt"
	"io"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

// TODO: race condition
type serverAPI struct {
	proto.UnimplementedOrchestratorServer
	tasks   chan *proto.Task
	results chan *proto.Task
	// mu sync.Mutex
}

func RegisterOrchestratorServer(gRPC *grpc.Server) {
	s := &serverAPI{
		tasks:   make(chan *proto.Task, 100),
		results: make(chan *proto.Task, 100),
	}
	go func() {
		for i := 0; i < 5; i++ {
			task := &proto.Task{ID: fmt.Sprint(i)}
			s.tasks <- task
		}
	}()
	proto.RegisterOrchestratorServer(gRPC, s)
}

func (s *serverAPI) TaskStream(stream grpc.BidiStreamingServer[proto.Task, proto.Task]) error {
	go func() {
		for task := range s.tasks {
			if err := stream.Send(task); err != nil {
				return
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			close(s.results)
			close(s.tasks)
			return nil
		}
		if err != nil {
			return err
		}
		s.results <- resp
		fmt.Println(resp)
	}
}
