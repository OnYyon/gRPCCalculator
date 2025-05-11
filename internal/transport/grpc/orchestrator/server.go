package orchestratorGRPC

import (
	"io"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

type serverAPI struct {
	manager *manager.Manager
	proto.UnimplementedOrchestratorServer
}

func RegisterOrchestratorServer(gRPC *grpc.Server, manager *manager.Manager) {
	s := &serverAPI{
		manager: manager,
	}
	proto.RegisterOrchestratorServer(gRPC, s)
}

func (s *serverAPI) TaskStream(stream grpc.BidiStreamingServer[proto.Task, proto.Task]) error {
	go func() {
		for task := range s.manager.Queque {
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
	}
}
