package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Executioner-OP/master/db"
	"github.com/Executioner-OP/master/pb"
	"github.com/Executioner-OP/master/queue"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExecutionsServer
}

var taskQueue = queue.Queue{}

func (s *server) GetExecution(context context.Context, executionRequest *pb.ExecutionRequest) (*pb.ExecutionTask, error) {
	var executionTask pb.ExecutionTask
	dummyExecutionTask := pb.ExecutionTask{
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
		HasTask:        false,
	}
	log.Println(taskQueue.GetLength())
	if taskQueue.IsEmpty() == false {
		task := taskQueue.Elements[0]
		executionTask.ID = task.ID.String()
		executionTask.Code = task.Code
		executionTask.ExpectedOutput = task.ExpectedOutput
		executionTask.StandardInput = task.StandardInput
		executionTask.StandardOutput = task.StandardOutput
		executionTask.IsDone = task.IsDone
		executionTask.Verdict = task.Verdict
		executionTask.Status = task.Status
		executionTask.TimeLimit = int32(task.TimeLimit)
		executionTask.MemoryLimit = int32(task.MemoryLimit)
		executionTask.LanguageId = int32(task.LanguageId)
		executionTask.HasTask = true
		taskQueue.Pop()
		return &executionTask, nil
	} else {
		return &dummyExecutionTask, nil
	}

}

func InitGrpcServer(taskChannel chan db.ExecutionRequest) {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen on Port 9001: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExecutionsServer(grpcServer, &server{})
	TaskChannel = taskChannel
	var wg sync.WaitGroup
	wg.Add(2)

	// Adding task from channel to queue concurrently
	go func() {
		for task := range TaskChannel {
			taskQueue.Add(task)
			log.Println("Task added to queue from channel")
		}
	}()

	// Start gRPC server
	go func() {
		defer wg.Done()
		fmt.Println("gRPC server running on port 9001")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server on port 9001: %v", err)
		}
	}()

	// Wait fot both servers to finish
	wg.Wait()

}
