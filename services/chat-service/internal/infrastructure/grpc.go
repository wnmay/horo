package infrastructure

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcin "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/grpc"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/proto/chat"
)

func SetupGRPCServer(app inbound_port.ChatService, addr string) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer()
	chatServer := grpcin.NewChatGRPCServer(app)
	chat.RegisterChatServiceServer(server, chatServer)

	reflection.Register(server)

	log.Printf("gRPC server listening on %s", addr)
	return server, lis, nil
}
