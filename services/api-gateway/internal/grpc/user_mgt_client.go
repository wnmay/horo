package grpcinfra

import (
	"context"
	"log"
	"time"

	pb "github.com/wnmay/horo/shared/proto/user_management"
	"google.golang.org/grpc"
)

type UserManagementClient struct {
	client pb.UserManagementServiceClient
}

func NewUserManagementClient(conn *grpc.ClientConn) *UserManagementClient {
	return &UserManagementClient{
		client: pb.NewUserManagementServiceClient(conn),
	}
}

func (oc *UserManagementClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := oc.client.Register(ctx, req)
	if err != nil {
		log.Printf("CreateUserManagement error: %v", err)
		return nil, err
	}

	return resp, nil
}
