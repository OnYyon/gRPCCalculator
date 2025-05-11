package api

import (
	"context"
	"fmt"
	"time"

	services "github.com/OnYyon/gRPCCalculator/internal/services/calculate"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (r *restAPI) Register(
	ctx context.Context,
	req *proto.AuthRequest,
) (*proto.AuthResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	// TODO: users exit
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	err = r.manager.DB.RegisterUser(context.TODO(), req.Login, hashedPassword)
	fmt.Println(req.Login, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthResponse{}, nil
}

func (r *restAPI) Login(
	ctx context.Context,
	req *proto.AuthRequest,
) (*proto.AuthResponse, error) {
	passwordHash, err := r.manager.DB.GetUser(context.TODO(), req.Login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Login,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(r.manager.Cfg.Auth.JWTSecret))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &proto.AuthResponse{Token: tokenString}, nil
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
	r.manager.AddStack(id, stack)
	if err != nil {
		return nil, err
	}
	err = r.manager.DB.SaveExpression(context.TODO(), id, request.Expression)
	if err != nil {
		panic(err)
	}
	return &proto.IDExpression{
		ID: id,
	}, nil
}
