package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/wnmay/horo/services/course-service/internal/domain"
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
func (c *UserClient) MapProphetNamesByIDs(ctx context.Context, userIDs []string) ([]domain.ProphetName, error) {
	req := &pb.MapProphetNamesRequest{
		UserIds: userIDs,
	}

	resp, err := c.client.MapProphetNames(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("nil response from user service")
	}

	log.Println("Successfully get prophet names mapping:", resp.Mappings)

	var prophetNames []domain.ProphetName
	for _, p := range resp.Mappings {
		prophetNames = append(prophetNames, domain.ProphetName{
			UserID: p.UserId,
			Name:   p.ProphetName,
		})
	}
	return prophetNames, nil
}

func (c *UserClient) GetProphetName(ctx context.Context, userID string) (string, error) {
	req := &pb.GetProphetNameRequest{
		UserId: userID,
	}
	resp, err := c.client.GetProphetName(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if resp == nil || resp.ProphetName == "" {
		return "", fmt.Errorf("empty response from user service")
	}

	prophetName := resp.ProphetName
	log.Println("Successfully get prophet name:", resp.ProphetName)

	return prophetName, nil
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
