package api

import (
	"context"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type restAPI struct {
	manager *manager.Manager
	proto.UnimplementedOrchestratorServer
}

func RegisterOrchestratorGateway(ctx context.Context, mux *runtime.ServeMux, manager *manager.Manager) error {
	s := &restAPI{
		manager: manager,
	}
	return proto.RegisterOrchestratorHandlerServer(ctx, mux, s)
}

func (r *restAPI) AddNewExpression(ctx context.Context, request *proto.Expression) (*proto.IDExpression, error) {
	return &proto.IDExpression{
		ID: "test",
	}, nil
}
