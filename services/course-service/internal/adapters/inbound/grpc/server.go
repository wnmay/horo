package grpc

import (
	"context"

	"github.com/wnmay/horo/services/course-service/internal/app"
	pb "github.com/wnmay/horo/shared/proto/course"
)

type CourseGRPCServer struct {
	pb.UnimplementedCourseServiceServer
	svc app.CourseService
}

func NewCourseGRPCServer(s app.CourseService) *CourseGRPCServer {
	return &CourseGRPCServer{svc: s}
}

func (s *CourseGRPCServer) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	in := app.CreateCourseInput{
		ProphetID:   req.GetProphetId(),
		CourseName:  req.GetCoursename(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Duration:    toDomainDuration(req.GetDuration()),
	}
	c, err := s.svc.CreateCourse(ctx, in)
	if err != nil {
		return nil, err
	}
	return &pb.CreateCourseResponse{Course: toPbCourse(c)}, nil
}

func (s *CourseGRPCServer) GetCourseByID(ctx context.Context, req *pb.GetCourseByIDRequest) (*pb.GetCourseByIDResponse, error) {
	c, err := s.svc.GetCourseByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetCourseByIDResponse{Course: toPbCourse(c)}, nil
}