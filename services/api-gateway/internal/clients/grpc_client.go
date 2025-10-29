package clients

import (
	grpc_connection "github.com/wnmay/horo/services/api-gateway/internal/services/grpc"
	pb "github.com/wnmay/horo/shared/proto/user-management"
)

type GrpcClients struct {
	AuthServiceClient pb.AuthServiceClient
	connManager       *grpc_connection.ConnectionManager
}

func NewGrpcClients(userManagementAddr string) (*GrpcClients, error) {
	cm := grpc_connection.NewConnectionManager()

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
