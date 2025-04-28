package orchestrator

import (
	"fmt"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

type serverAPI struct {
	// Its plug for test
	proto.UnimplementedOrchestratorServer
	Transport
}

type Transport interface {
	TransportTasks(proto.Task) proto.Task
}

func RegisterOrchestratorServer(gRPC *grpc.Server) {
	proto.RegisterOrchestratorServer(gRPC, &serverAPI{})
}

// TODO: доделать логику
func (s *serverAPI) TransportTasks(
	stream grpc.BidiStreamingServer[proto.Task, proto.Task],
) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	fmt.Println(req)
	return nil
}
