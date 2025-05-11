package main

import (
	"fmt"
	"log"

	"github.com/Retual004/tasks-service/internal/database"
	"github.com/Retual004/tasks-service/internal/task"
	transportgrpc "github.com/Retual004/tasks-service/internal/transport/grpc"
)
func main() {
	addr:= ":50052"
	database.InitDB()

	repo:= task.NewRepository(database.DB)
	svc:= task.NewService(repo)
 	
	userClient, conn, err := transportgrpc.NewUserClient("localhost:50051")
  if err != nil {
    log.Fatalf("failed to connect to users: %v", err)
  }
  defer conn.Close()

  fmt.Printf("gRPC server listening on %s\n", addr)
    // 4. Запуск gRPC Tasks-сервиса
  if err := transportgrpc.RunGRPC(svc, userClient); err != nil {
    log.Fatalf("Tasks gRPC server error: %v", err)
  }
}
