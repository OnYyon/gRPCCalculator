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
	host := "localhost"
	port := "8080"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части GeometryService
	// зарегистрируем нашу реализацию сервера
	orchestrator.RegisterOrchestratorServer(grpcServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
