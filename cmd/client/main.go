package main

import (
	"context"
	"fmt"
	"io"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type quequeTask struct {
	queque chan *proto.Task
}

// For tests
// Я думаю жто кастыль, но он работает
// TODO: Попробыть способ Unary и по изучать еще как это реализоватья
func main() {
	numWorkers := 3
	quequeTask := quequeTask{make(chan *proto.Task, numWorkers)}
	for i := 0; i < numWorkers; i++ {
		go quequeTask.StartAgent(i)
	}
	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("erorr!")
	}
	defer conn.Close()
	grpcClient := proto.NewOrchestratorClient(conn)
	stream, err := grpcClient.TransportTasks(context.Background())
	if err != nil {
		panic("error in client/main.go")
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		go func() { quequeTask.queque <- req }()
	}
}

func (q *quequeTask) StartAgent(num int) {
	task := <-q.queque
	fmt.Println(num, task)

}
