package grpcadapter

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/ports"
	proto "github.com/wnmay/horo/shared/proto/user-management"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	authService ports.AuthService
}

func NewAuthServer(authService ports.AuthService) *AuthServer {
	return &AuthServer{authService: authService}
}

func (s *AuthServer) GetClaims(ctx context.Context, req *proto.GetClaimsRequest) (*proto.GetClaimsResponse, error) {
	claims, err := s.authService.GetClaims(ctx, req.IdToken)
	if err != nil {
		return nil, err
	}
	return &proto.GetClaimsResponse{
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}
