package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/OnYyon/gRPCCalculator/iternal/grpc/orchestrator"
	"google.golang.org/grpc"
)

func main() {
	// TODO: загрузка конфигуряция из .env

	host := "localhost"
	port := "8080"

	address := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic("dont start listen")
	}
	fmt.Println("starting listen on :8080")

	grpcServer := grpc.NewServer()
	orchestrator.RegisterOrchestratorServer(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
