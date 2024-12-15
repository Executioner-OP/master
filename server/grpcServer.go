package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Executioner-OP/master/db"
	"github.com/Executioner-OP/master/pb"
	"github.com/Executioner-OP/master/queue"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExecutionsServer
}

var (
	taskQueue   = queue.Queue{}
	taskTimeout = 10 * time.Second
)

func (s *server) GetExecution(ctx context.Context, req *pb.ExecutionRequest) (*pb.ExecutionTask, error) {
	if taskQueue.IsEmpty() {
		return &pb.ExecutionTask{HasTask: false}, nil
	}
	task, _ := taskQueue.Pop()

	return &pb.ExecutionTask{
		ID:             task.ID.String(),
		Code:           task.Code,
		ExpectedOutput: task.ExpectedOutput,
		StandardInput:  task.StandardInput,
		TimeLimit:      int32(task.TimeLimit),
		MemoryLimit:    int32(task.MemoryLimit),
		LanguageId:     int32(task.LanguageId),
		Status:         "pending",
		HasTask:        true,
	}, nil
}

func InitGrpcServer(taskChannel chan db.ExecutionRequest) {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen on port 9001: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExecutionsServer(grpcServer, &server{})

	go func() {
		for task := range taskChannel {
			taskQueue.Add(task, taskTimeout)
		}
	}()

	fmt.Println("gRPC server running on port 9001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
