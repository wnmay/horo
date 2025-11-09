package grpcadapter

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
	proto "github.com/wnmay/horo/shared/proto/user-management"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	userManagementService ports.UserManagementService
}

func NewUserServer(userService ports.UserManagementService) *UserServer {
	return &UserServer{userManagementService: userService}
}

func (s *UserServer) MapProphetNames(ctx context.Context, req *proto.MapProphetNamesRequest) (*proto.MapProphetNamesResponse, error) {
	mappings, err := s.userManagementService.GetProphetNames(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	protoMappings := make([]*proto.ProphetData, len(mappings))
	for i, mapping := range mappings {
		protoMappings[i] = toProtoProphetName(mapping)
	}
	return &proto.MapProphetNamesResponse{
		Mappings: protoMappings,
	}, nil
}

func (s *UserServer) GetProphetName(ctx context.Context, req *proto.GetProphetNameRequest) (*proto.GetProphetNameResponse, error) {
	prophetName, err := s.userManagementService.GetProphetName(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.GetProphetNameResponse{
		ProphetName: prophetName,
	}, nil
}

func toProtoProphetName(prophetName *domain.ProphetName) *proto.ProphetData {
	return &proto.ProphetData{
		UserId:      prophetName.UserID,
		ProphetName: prophetName.ProphetName,
	}
}
