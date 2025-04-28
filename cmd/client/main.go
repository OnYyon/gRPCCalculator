package main

import (
	"context"
	"fmt"
	"log"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// For tests
func main() {
	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("error!")
	}
	defer conn.Close()

	grpcClient := proto.NewOrchestratorClient(conn)
	stream, err := grpcClient.TransportTasks(context.Background())
	if err != nil {
		log.Fatalf("failed to create stream: %v", err)
	}
	task, err := stream.Recv()
	if err != nil {
		panic("error!")
	}
	fmt.Printf("get: %v", task)
}
