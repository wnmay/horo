package clients

import (
	grpc_connection "github.com/wnmay/horo/services/api-gateway/internal/services/grpc"
	chatpb "github.com/wnmay/horo/shared/proto/chat"
	authpb "github.com/wnmay/horo/shared/proto/user-management"
)

type GrpcClients struct {
	AuthServiceClient authpb.AuthServiceClient
	connManager       *grpc_connection.ConnectionManager
	ChatServiceClient chatpb.ChatServiceClient
}

func NewGrpcClients(userManagementAddr string, chatServiceAddr string) (*GrpcClients, error) {
	cm := grpc_connection.NewConnectionManager()

	userManagementConn, err := cm.GetConnection(userManagementAddr)
	if err != nil {
		return nil, err
	}
	authClient := authpb.NewAuthServiceClient(userManagementConn)

	chatConn, err := cm.GetConnection(chatServiceAddr)
	if err != nil {
		return nil, err
	}
	chatClient := chatpb.NewChatServiceClient(chatConn)

	return &GrpcClients{
		AuthServiceClient: authClient,
		connManager:       cm,
		ChatServiceClient: chatClient,
	}, nil
}

func (gc *GrpcClients) Close() {
	gc.connManager.CloseAll()
}
