package service

func GeneratePaymentCreatedMessage(paymentID string, orderID string, status string, amount float64) string {
	return `
		Payment created successfully for order %s with payment ID %s. Status: %s, Amount: %f
	`
}

func GenerateOrderCompletedMessage(orderID string, courseID string, orderStatus string, courseName string) string {
	return `
		Order %s completed successfully for course %s. Status: %s, Course Name: %s
	`
}

func GenerateOrderPaymentBoundMessage(orderID string, courseID string, orderStatus string, courseName string) string {
	return `
		Successfully create payment for order %s payment created successfully for course %s. Status: %s, Course Name: %s
	`
}

func GenerateOrderPaidMessage(orderID string, courseID string, orderStatus string, courseName string) string {
	return `
		Order %s paid successfully for course %s. Status: %s, Course Name: %s
	`
}