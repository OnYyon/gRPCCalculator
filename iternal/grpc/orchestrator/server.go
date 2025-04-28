package orchestrator

import (
	"context"
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

func (s *serverAPI) AddNewExpression(ctx context.Context, in *proto.ExpressionID) (*proto.ResponseID, error) {
	fmt.Println(in)
	return nil, nil
}

func (s *serverAPI) GetExpressionByID(ctx context.Context, in *proto.ExpressionID) (*proto.Expression, error) {
	fmt.Println(in)
	return nil, nil
}

func (s *serverAPI) GetAllExpressions(ctx context.Context, in *proto.Empty) (*proto.ExpressionsList, error) {
	fmt.Println(in)
	return nil, nil
}
