package grpcinfra

import (
	pb "github.com/wnmay/horo/shared/proto/user-management"
)

type GrpcClients struct {
	AuthServiceClient pb.AuthServiceClient
	connManager       *ConnectionManager
}

func NewGrpcClients(userManagementAddr string) (*GrpcClients, error) {
	cm := NewConnectionManager()

	userManagementConn, err := cm.GetConnection(userManagementAddr)
	if err != nil {
		return nil, err
	}
	authClient := pb.NewAuthServiceClient(userManagementConn)

	return &GrpcClients{
		AuthServiceClient: authClient,
		connManager:       cm,
	}, nil
}

func (gc *GrpcClients) Close() {
	gc.connManager.CloseAll()
}
