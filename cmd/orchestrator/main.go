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
	// TODO: сделать логирование

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
	// for i := 0; i < 102; i++ {
	// 	s.AddTask(&proto.Task{ID: fmt.Sprint(i), Arg1: 1.0, Arg2: 2.0, Operator: "+", Have: true})
	// }
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
