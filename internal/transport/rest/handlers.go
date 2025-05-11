package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	services "github.com/OnYyon/gRPCCalculator/internal/services/calculate"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	"github.com/OnYyon/gRPCCalculator/internal/transport/grpc/auth"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type restAPI struct {
	manager *manager.Manager
	auth    *auth.AuthGRPC
	proto.UnimplementedOrchestratorServer
}

func RegisterOrchestratorGateway(
	ctx context.Context,
	mux *runtime.ServeMux,
	manager *manager.Manager,
	auth *auth.AuthGRPC,
) error {
	s := &restAPI{
		manager: manager,
		auth:    auth,
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

// TODO: улучшить струткру.
func (r *restAPI) AddNewExpression(
	ctx context.Context,
	request *proto.Expression,
) (*proto.IDExpression, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
	}
	token := strings.TrimPrefix(authHeaders[0], "Bearer ")

	userid, err := r.auth.ValidateTokenAndGetUserID(token)
	if err != nil {
		return nil, err
	}
	id := r.manager.GenerateUUID()
	rpn, err := services.ParserToRPN(request.Expression)
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
	err = r.manager.DB.SaveExpression(context.TODO(), id, request.Expression, userid)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	return &proto.IDExpression{
		Id: id,
	}, nil
}

func (r *restAPI) GetExpressionByID(
	ctx context.Context,
	request *proto.IDExpression,
) (*proto.ExpressionRes, error) {
	m, err := r.manager.DB.GetExpressionByID(context.TODO(), request.Id)
	if err != nil {
		return nil, err
	}
	o := &proto.ExpressionRes{
		ID:     m["id"],
		Status: m["status"],
		Result: m["result"],
		Input:  m["expression"],
	}
	return o, nil
}

func (r *restAPI) GetListExpression(
	ctx context.Context,
	requset *proto.TNIL,
) (*proto.ExpressionList, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
	}
	token := strings.TrimPrefix(authHeaders[0], "Bearer ")

	userid, err := r.auth.ValidateTokenAndGetUserID(token)
	if err != nil {
		return nil, err
	}
	slc, err := r.manager.DB.GetExpressionList(context.TODO(), userid)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	var expressionList proto.ExpressionList
	for _, expr := range slc {
		expressionRes := &proto.ExpressionRes{
			ID:     expr.ID,
			Status: expr.Status,
			Result: expr.Result,
			Input:  expr.Expression,
		}
		expressionList.List = append(expressionList.List, expressionRes)
	}
	return &expressionList, nil
}
