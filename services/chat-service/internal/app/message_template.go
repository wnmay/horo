package service

func GeneratePaymentCreatedMessage(paymentID string, orderID string, status string, amount float64) string {
	return `
	<div class="message-container">
		<div class="message-header">
			<h3>Payment Created</h3>
		</div>
		<div class="message-body">
			<p>Payment created successfully for order %s with payment ID %s. Status: %s, Amount: %f</p>
		</div>
	</div>
	`
}
