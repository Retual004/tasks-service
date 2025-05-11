package grpc

import (
	"github.com/Retual004/project-protos/proto/user"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc"
)

func NewUserClient(addr string) (user.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		addr,  
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client := user.NewUserServiceClient(conn)
	    return client, conn, nil
}