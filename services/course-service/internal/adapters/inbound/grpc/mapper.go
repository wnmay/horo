package grpc

import (
	"github.com/wnmay/horo/services/course-service/internal/domain"
	pb "github.com/wnmay/horo/shared/proto/course"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPbDuration(d domain.DurationEnum) pb.Duration {
	switch d {
	case 15:
		return pb.Duration_DURATION_15
	case 30:
		return pb.Duration_DURATION_30
	case 60:
		return pb.Duration_DURATION_60
	default:
		return pb.Duration_DURATION_UNSPECIFIED
	}
}

func toDomainDuration(d pb.Duration) domain.DurationEnum {
	switch d {
	case pb.Duration_DURATION_15:
		return 15
	case pb.Duration_DURATION_30:
		return 30
	case pb.Duration_DURATION_60:
		return 60
	default:
		return 0
	}
}

func toPbCourse(c *domain.Course) *pb.Course {
	if c == nil {
		return nil
	}
	return &pb.Course{
		Id:          c.ID,
		ProphetId:   c.ProphetID,
		Coursename:  c.CourseName,
		Description: c.Description,
		Price:       c.Price,
		Duration:    toPbDuration(c.Duration),
		CreatedTime: timestamppb.New(c.CreatedAt),
	}
}
