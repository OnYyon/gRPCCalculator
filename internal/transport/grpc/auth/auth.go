package auth

import (
	"context"
	"fmt"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthGRPC struct {
	manager *manager.Manager
}

func NewAuthGRPC(mgr *manager.Manager) *AuthGRPC {
	return &AuthGRPC{manager: mgr}
}

func (a *AuthGRPC) AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info.FullMethod == "/orchestrator.Orchestrator/Register" ||
		info.FullMethod == "/orchestrator.Orchestrator/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	tokenString := authHeader[0]
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	tokenString = tokenString[7:]

	login, err := a.ValidateTokenAndGetUserID(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	_, err = a.manager.DB.GetUser(context.TODO(), login)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}
	return handler(ctx, req)
}

func (a *AuthGRPC) ValidateTokenAndGetUserID(tokenString string) (string, error) {

	if a.manager.Cfg.Auth.JWTSecret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}

	secretKey := []byte(a.manager.Cfg.Auth.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("invalid user ID in token")
		}
		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}
