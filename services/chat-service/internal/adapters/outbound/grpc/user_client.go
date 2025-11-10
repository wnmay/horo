package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	pb "github.com/wnmay/horo/shared/proto/user-management"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient creates a new user service gRPC client
func NewUserClient(userServiceAddr string) (*UserClient, error) {
	conn, err := grpc.NewClient(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := pb.NewUserServiceClient(conn)
	log.Printf("Connected to user service at %s", userServiceAddr)

	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// GetUserByID fetches user details by ID
func (c *UserClient) MapUserNamesByIDs(ctx context.Context, userIDs []string) (map[string]*domain.User, error) {
	req := &pb.MapUserNamesRequest{
		UserIds: userIDs,
	}

	resp, err := c.client.MapUserNames(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("nil response from user service")
	}

	log.Println("Successfully got user name mapping:", resp.Users)

	users := make(map[string]*domain.User, len(resp.Users))
	for id, user := range resp.Users {
		users[id] = &domain.User{
			Name: user.Name,
			Role: user.Role.String(),
		}
	}

	return users, nil
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
