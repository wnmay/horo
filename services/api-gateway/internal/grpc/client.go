package grpcinfra

import (
	pb "github.com/wnmay/horo/shared/proto/user-management"
)

type GrpcClients struct {
	UserManagementClient pb.UserManagementServiceClient
	connManager          *ConnectionManager
}

func NewGrpcClients(userManagementAddr string) (*GrpcClients, error) {
	cm := NewConnectionManager()

	userManagementConn, err := cm.GetConnection(userManagementAddr)
	if err != nil {
		return nil, err
	}
	userManagementClient := pb.NewUserManagementServiceClient(userManagementConn)

	return &GrpcClients{
		UserManagementClient: userManagementClient,
		connManager:          cm,
	}, nil
}

func (gc *GrpcClients) Close() {
	gc.connManager.CloseAll()
}
