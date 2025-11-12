package service

import "fmt"

func GeneratePaymentCreatedMessage(paymentID string, orderID string, status string, amount float64) string {
	return fmt.Sprintf("Payment created successfully for order %s with payment ID %s. Status: %s, Amount: %.2f",
		orderID, paymentID, status, amount)
}

func GenerateOrderCompletedMessage(orderID string, courseID string, orderStatus string, courseName string) string {
	return fmt.Sprintf("Order %s completed successfully for course %s. Status: %s, Course Name: %s",
		orderID, courseID, orderStatus, courseName)
}

func GenerateOrderPaymentBoundMessage(orderID string, courseID string, orderStatus string, courseName string) string {
	return fmt.Sprintf("Successfully create payment for order %s payment created successfully for course %s. Status: %s, Course Name: %s",
		orderID, courseID, orderStatus, courseName)
}

func GenerateOrderPaidMessage(orderID string, courseID string, orderStatus string, courseName string) string {
	return fmt.Sprintf("Order %s paid successfully for course %s. Status: %s, Course Name: %s",
		orderID, courseID, orderStatus, courseName)
}
