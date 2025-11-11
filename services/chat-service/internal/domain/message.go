package domain

import (
	"time"
)

type Message struct {
	ID        string
	RoomID    string
	SenderID  string
	Content   string
	Type      MessageType   // text | notification
	Trigger   string        // order.created | order.completed | order.payment.bound | order.paid | payment.completed | payment.created | payment.settled
	Status    MessageStatus // sent | delivered | read
	CreatedAt time.Time
}

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

type MessageType string

const (
	MessageTypeText         MessageType = "text"
	MessageTypeNotification MessageType = "notification"
)

func CreateMessage(messageID, roomID, senderID, content string, messageType MessageType, status MessageStatus, trigger string) *Message {
	return &Message{
		ID:        messageID,
		RoomID:    roomID,
		SenderID:  senderID,
		Content:   content,
		Type:      messageType,
		Status:    status,
		CreatedAt: time.Now(),
		Trigger:   trigger,
	}
}
