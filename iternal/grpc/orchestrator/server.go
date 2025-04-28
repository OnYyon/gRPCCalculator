package orchestrator

import (
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
)

type serverAPI struct {
	proto.UnimplementedOrchestratorServer
	queque chan *proto.Task
}

func RegisterOrchestratorServer(gRPC *grpc.Server) {
	proto.RegisterOrchestratorServer(gRPC, &serverAPI{})
}

// TODO: доделать логику
func (s *serverAPI) TransportTasks(
	stream grpc.BidiStreamingServer[proto.Task, proto.Task],
) error {
	sliceTests := []*proto.Task{
		{
			ID:           "1",
			Arg1:         "2",
			Arg2:         "2",
			Operation:    "+",
			ExpressionID: "1",
		},
		{
			ID:           "2",
			Arg1:         "3",
			Arg2:         "3",
			Operation:    "+",
			ExpressionID: "2",
		},
		{
			ID:           "3",
			Arg1:         "4",
			Arg2:         "4",
			Operation:    "+",
			ExpressionID: "3",
		},
	}

	for _, task := range sliceTests {
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}
