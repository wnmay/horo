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
}

func CreateRoom(prophetID, customerID, courseID string) *Room {
	return &Room{
		ProphetID:  prophetID,
		CustomerID: customerID,
		CourseID:   courseID,
		CreatedAt:  time.Now(),
	}
}
