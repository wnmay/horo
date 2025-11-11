package outbound_port

import (
	"context"
)

type CourseProvider interface {
	GetProphetIDByCourseID(ctx context.Context, courseID string) (string, error)
}
