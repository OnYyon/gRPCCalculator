package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	orchestratorGRPC "github.com/OnYyon/gRPCCalculator/internal/transport/grpc/orchestrator"
	api "github.com/OnYyon/gRPCCalculator/internal/transport/rest"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type App struct {
	cfg     *config.Config
	manager *manager.Manager
}

func New(cfg *config.Config) *App {
	mgr := manager.NewManager(cfg)
	return &App{
		cfg:     cfg,
		manager: mgr,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lis, err := net.Listen("tcp", a.cfg.Server.Host+":"+a.cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	m := cmux.New(lis)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	go a.runGRPCServer(grpcL)
	go a.runHTTPServer(ctx, httpL)

	log.Printf("Server started on %s:%s", a.cfg.Server.Host, a.cfg.Server.Port)
	return m.Serve()
}

func (a *App) runGRPCServer(lis net.Listener) {
	grpcServer := grpc.NewServer()
	orchestratorGRPC.RegisterOrchestratorServer(grpcServer, a.manager)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("GRPC server failed: %v", err)
		os.Exit(1)
	}
}

func (a *App) runHTTPServer(ctx context.Context, lis net.Listener) {
	mux := runtime.NewServeMux()

	if err := api.RegisterOrchestratorGateway(ctx, mux, a.manager); err != nil {
		log.Printf("Failed to register gateway: %v", err)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on", a.cfg.Server.Host+":"+a.cfg.Server.Port)
	if err := httpServer.Serve(lis); err != nil && err != cmux.ErrListenerClosed {
		log.Printf("HTTP server failed: %v", err)
		os.Exit(1)
	}
}

func (a *App) Close() error {
	if a.manager.DB != nil {
		return a.manager.DB.Close()
	}
	return nil
}
