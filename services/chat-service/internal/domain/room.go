package domain

import (
	"time"
)

type Room struct {
	ID          string
	ProphetID   string
	CustomerID  string
	CourseID    string
	CreatedAt   time.Time
	LastMessage string
	IsDone      bool
}

func CreateRoom(prophetID, customerID, courseID string, isDone bool) *Room {
	return &Room{
		ProphetID:  prophetID,
		CustomerID: customerID,
		CourseID:   courseID,
		CreatedAt:  time.Now(),
		IsDone:     isDone,
	}
}
