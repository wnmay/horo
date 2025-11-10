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

func (s *UserServer) GetProphetIdsByNames(ctx context.Context, req *proto.GetProphetIdsByNamesRequest) (*proto.GetProphetIdsByNamesResponse, error) {
	prophetIds, err := s.userManagementService.SearchProphetIdsByName(ctx, req.ProphetName)
	if err != nil {
		return nil, err
	}
	protoProphetData := make([]*proto.ProphetData, len(prophetIds))
	for i, prophetId := range prophetIds {
		protoProphetData[i] = toProtoProphetName(prophetId)
	}
	return &proto.GetProphetIdsByNamesResponse{
		ProphetData: protoProphetData,
	}, nil
}

func (s *UserServer) MapUserNames(ctx context.Context, req *proto.MapUserNamesRequest) (*proto.MapUserNamesResponse, error) {
	userNames, err := s.userManagementService.MapUserNames(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	protoUserNames := make([]*proto.UserData, len(userNames))
	for i, userName := range userNames {
		protoUserNames[i] = toProtoUserName(userName)
	}
	return &proto.MapUserNamesResponse{
		Users: protoUserNames,
	}, nil
}

func toProtoProphetName(prophetName *domain.ProphetName) *proto.ProphetData {
	return &proto.ProphetData{
		UserId:      prophetName.UserID,
		ProphetName: prophetName.ProphetName,
	}
}

func toProtoUserName(userName *domain.UserName) *proto.UserData {
	return &proto.UserData{
		UserId: userName.UserID,
		Name:   userName.UserName,
		Role:   proto.UserRole(proto.UserRole_value[string(userName.UserRole)]),
	}
}
