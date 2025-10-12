// internal/adapter/grpc/handler.go
package grpcadapter

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/ports"
	proto "github.com/wnmay/horo/shared/proto/user-management"
)

type GRPCServer struct {
	proto.UnimplementedUserManagementServiceServer
	service ports.UserManagementService
}

func NewGRPCServer(s ports.UserManagementService) *GRPCServer {
	return &GRPCServer{service: s}
}

func (g *GRPCServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	err := g.service.Register(ctx, req.FirebaseToken, req.FullName, req.Role)
	if err != nil {
		return &proto.RegisterResponse{
			Success: false,
		}, err
	}
	return &proto.RegisterResponse{
		Success: true,
	}, nil
}

func (g *GRPCServer) GetClaims(ctx context.Context, req *proto.GetClaimsRequest) (*proto.GetClaimsResponse, error) {
	claims, err := g.service.ValidateFirebaseToken(req.FirebaseToken)
	if err != nil {
		return nil, err
	}
	return &proto.GetClaimsResponse{
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}
