package orchestrator

import (
	"fmt"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

type serverAPI struct {
	proto.UnimplementedOrchestratorServer
}

func RegisterOrchestratorServer(gRPC *grpc.Server) {
	proto.RegisterOrchestratorServer(gRPC, &serverAPI{})
}

// TODO: доделать логику
func (s *serverAPI) TransportTasks(
	stream grpc.BidiStreamingServer[proto.Task, proto.Task],
) error {
	task := &proto.Task{
		ID:           "1",
		Arg1:         "2",
		Arg2:         "2",
		Operation:    "+",
		ExpressionID: "1",
	}
	if err := stream.Send(task); err != nil {
		panic("error!")
	}
	fmt.Printf("recieve %v", task)
	return nil
}
