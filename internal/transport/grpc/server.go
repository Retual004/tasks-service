package grpc


import(
	"github.com/Retual004/tasks-service/internal/task"
	userpb "github.com/Retual004/project-protos/proto/user"
	taskpb "github.com/Retual004/project-protos/proto/task"
	"google.golang.org/grpc"
	"net"
)

func RunGRPC(svc *task.TaskService, uc userpb.UserServiceClient) error {
  lis, _ := net.Listen("tcp", ":50052")
  grpcSrv := grpc.NewServer()
  handler := NewHandler(svc, uc)
  taskpb.RegisterTaskServiceServer(grpcSrv, handler)
  return grpcSrv.Serve(lis)
}