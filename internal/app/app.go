package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	orchestratorGRPC "github.com/OnYyon/gRPCCalculator/internal/transport/grpc/orchestrator"
	api "github.com/OnYyon/gRPCCalculator/internal/transport/rest"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

func rungRPC(lis net.Listener) {
	// TODO: сделать чтение из .env
	manager := manager.NewManager()
	grpcServer := grpc.NewServer()
	orchestratorGRPC.RegisterOrchestratorServer(grpcServer, manager)
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

func runRestAPI(lis net.Listener) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	manager := manager.NewManager()
	if err := api.RegisterOrchestratorGateway(ctx, mux, manager); err != nil {
		panic(err)
	}
	httpServer := &http.Server{
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on :8080")
	if err := httpServer.Serve(lis); err != cmux.ErrListenerClosed {
		log.Fatalf("failed to serve gateway: %v", err)
	}
}

func StartOrchestrator(cfg *config.Config) {
	port := cfg.Server.Port
	lis, err := net.Listen("tcp", cfg.Server.Host+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m := cmux.New(lis)

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	// Start servers
	go rungRPC(grpcL)
	go runRestAPI(httpL)

	log.Println("Server started on", cfg.Server.Host+":"+port)
	if err := m.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
