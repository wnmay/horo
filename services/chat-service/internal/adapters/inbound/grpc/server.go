// internal/adapters/in/grpc/server.go
package grpcin

import (
	"context"

	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/proto/chat"
)

type ChatGRPCServer struct {
	chat.UnimplementedChatServiceServer
	app inbound_port.ChatService
}

func NewChatGRPCServer(app inbound_port.ChatService) *ChatGRPCServer {
	return &ChatGRPCServer{app: app}
}

func (s *ChatGRPCServer) ValidateRoomAccess(ctx context.Context, req *chat.ValidateRoomRequest) (*chat.ValidateRoomResponse, error) {
	allow, reason, err := s.app.ValidateRoomAccess(ctx, req.UserId,req.RoomId)
	if err != nil {
		return nil, err
	}
	return &chat.ValidateRoomResponse{
		Allowed: allow,
		Reason:  reason,
	}, nil
}