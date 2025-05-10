package api

import (
	"context"
	"fmt"
	"time"

	services "github.com/OnYyon/gRPCCalculator/internal/services/calculate"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type restAPI struct {
	manager *manager.Manager
	proto.UnimplementedOrchestratorServer
}

func RegisterOrchestratorGateway(
	ctx context.Context,
	mux *runtime.ServeMux,
	manager *manager.Manager,
) error {
	s := &restAPI{
		manager: manager,
	}
	return proto.RegisterOrchestratorHandlerServer(ctx, mux, s)
}

// TODO: доделать полный цикл решения задачи.
// TODO: улучшить струткру.
func (r *restAPI) AddNewExpression(
	ctx context.Context,
	request *proto.Expression,
) (*proto.IDExpression, error) {
	id := r.manager.GenerateUUID()
	rpn, err := services.ParserToRPN(request.Expression)
	fmt.Println(rpn)
	if err != nil {
		return nil, err
	}
	stack, tasks, err := services.GenerateTasks(rpn, id, r.manager)
	for _, task := range tasks {
		r.manager.AddTask(task)
	}
	time.Sleep(3 * time.Second)
	fmt.Println(stack, r.manager.Expressions)
	if err != nil {
		return nil, err
	}
	err = r.manager.DB.SaveExpression(context.TODO(), id, id)
	if err != nil {
		panic(err)
	}
	return &proto.IDExpression{
		ID: id,
	}, nil
}
