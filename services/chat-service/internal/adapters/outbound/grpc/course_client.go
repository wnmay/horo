package grpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/wnmay/horo/shared/proto/course"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CourseClient struct {
	client pb.CourseServiceClient
	conn   *grpc.ClientConn
}

// NewCourseClient creates a new course service gRPC client
func NewCourseClient(courseServiceAddr string) (*CourseClient, error) {
	conn, err := grpc.NewClient(
		courseServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to course service: %w", err)
	}

	client := pb.NewCourseServiceClient(conn)
	log.Printf("✅ Connected to course service at %s", courseServiceAddr)

	return &CourseClient{
		client: client,
		conn:   conn,
	}, nil
}

// GetCourseByID fetches course details by ID
func (c *CourseClient) GetCourseByID(ctx context.Context, courseID string) (*pb.Course, error) {
	req := &pb.GetCourseByIDRequest{
		Id: courseID,
	}

	log.Printf("Calling course service GetCourseByID for courseID: %s", courseID)

	resp, err := c.client.GetCourseByID(ctx, req)
	if err != nil {
		log.Printf("gRPC error calling GetCourseByID: %v", err)
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	if resp == nil {
		log.Printf("GetCourseByID returned nil response")
		return nil, fmt.Errorf("nil response from course service")
	}

	if resp.Course == nil {
		log.Printf("GetCourseByID response has nil course")
		return nil, fmt.Errorf("course not found")
	}

	log.Printf("Successfully fetched course: ID=%s, Name=%s, Price=%.2f",
		resp.Course.Id, resp.Course.Coursename, resp.Course.Price)

	return resp.Course, nil
}

// GetProphetIDByCourseID fetches the prophet ID for a given course
func (c *CourseClient) GetProphetIDByCourseID(ctx context.Context, courseID string) (string, error) {
	course, err := c.GetCourseByID(ctx, courseID)
	if err != nil {
		return "", fmt.Errorf("failed to get course for prophet ID: %w", err)
	}

	if course.ProphetId == "" {
		return "", fmt.Errorf("course %s has no prophet ID", courseID)
	}

	log.Printf("Found prophet ID %s for course %s", course.ProphetId, courseID)
	return course.ProphetId, nil
}

// Close closes the gRPC connection
func (c *CourseClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
