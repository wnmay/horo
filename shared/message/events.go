package message

const (
	CreatePaymentQueue = "create_payment_queue"
	UpdateOrderStatusQueue = "update_order_status_queue"

	DeadLetterQueue = "dead_letter_queue"
)


// ---- DATA STRUCTURES ----

type OrderData struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
}

type PaymentSuccessData struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	TransactionID string `json:"transaction_id"`
}