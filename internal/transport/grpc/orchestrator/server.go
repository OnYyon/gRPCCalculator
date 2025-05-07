package orchestratorGRPC

import (
	"fmt"
	"io"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

// TODO: race condition
type serverAPI struct {
	manager *manager.Manager
	proto.UnimplementedOrchestratorServer
	// mu sync.Mutex
}

func RegisterOrchestratorServer(gRPC *grpc.Server, manager *manager.Manager) {
	s := &serverAPI{
		manager: manager,
	}
	// NOTE: for tests
	go func() {
		for i := 0; i < 5; i++ {
			task := &proto.Task{ID: fmt.Sprint(i)}
			s.manager.AddTask(task)
		}
	}()
	proto.RegisterOrchestratorServer(gRPC, s)
}

func (s *serverAPI) TaskStream(stream grpc.BidiStreamingServer[proto.Task, proto.Task]) error {
	go func() {
		for task := range s.manager.Tasks {
			if err := stream.Send(task); err != nil {
				return
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		s.manager.AddResult(resp)
		fmt.Println(resp)
	}
}
