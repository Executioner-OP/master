package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Executioner-OP/master/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExecutionsServer
}

func (s *server) GetExecution(context context.Context, executionRequest *pb.ExecutionRequest) (*pb.ExecutionTask, error) {
	var dummyExecutionTask pb.ExecutionTask
	dummyExecutionTask = pb.ExecutionTask{
		ID:             "60d5ec49f1d2c12a4c8b4567",
		Code:           "#include <iostream>\nint main() {\n    std::cout << \"Hello, World!\" << std::endl;\n    return 0;\n}",
		IsDone:         false,
		LanguageId:     1,
		StandardInput:  "input data",
		StandardOutput: "",
		ExpectedOutput: "Hello, World!",
		Status:         "pending",
		Verdict:        "",
		TimeLimit:      5,
		MemoryLimit:    256,
	}

	return &dummyExecutionTask, nil
}

func InitGrpcServer() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen on Port 9001: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterExecutionsServer(grpcServer, &server{})
	fmt.Println("gRPC server running on port 9001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 9001: %v", err)
	}

}
