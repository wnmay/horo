package grpcadapter

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/ports"
	proto "github.com/wnmay/horo/shared/proto/user-management"
)

type UserManagementServer struct {
	proto.UnimplementedUserManagementServiceServer
	userService ports.UserManagementService
}

func NewUserManagementServer(userService ports.UserManagementService) *UserManagementServer {
	return &UserManagementServer{userService: userService}
}

func (s *UserManagementServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	err := s.userService.Register(ctx, req.FirebaseToken, req.FullName, req.Role)
	if err != nil {
		return &proto.RegisterResponse{Success: false}, err
	}
	return &proto.RegisterResponse{Success: true}, nil
}
