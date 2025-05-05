package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	orchestratorGRPC "github.com/OnYyon/gRPCCalculator/internal/transport/grpc/orchestrator"
	"google.golang.org/grpc"
)

func main() {
	// TODO: загрузка конфигуряция из .env
	// TODO: сделать логирование

	host := "localhost"
	port := "8080"

	address := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic("dont start listen")
	}
	fmt.Println("starting listen on :8080")
	manager := manager.NewManager()
	grpcServer := grpc.NewServer()
	orchestratorGRPC.RegisterOrchestratorServer(grpcServer, manager)
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
