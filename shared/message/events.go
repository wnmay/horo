package message

const (
	CreatePaymentQueue       = "create_payment_queue"
	UpdateOrderStatusQueue   = "update_order_status_queue"
	UpdatePaymentIDQueue     = "update_payment_id_queue"
	ChatMessageIncomingQueue = "chat_message_incoming_queue"
	ChatMessageOutgoingQueue = "chat_message_outgoing_queue"
	DeadLetterQueue          = "dead_letter_queue"
	SettlePaymentQueue       = "settle_payment_queue"
	NotifyCreatePayment      = "notify_create_payment"
)

// ---- DATA STRUCTURES ----

type OrderData struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
}

type PaymentSuccessData struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	TransactionID string `json:"transaction_id"`
}

type ChatMessageIncomingData struct {
	RoomID   string `json:"roomId"`
	SenderID string `json:"senderId"`
	Content  string `json:"content"`
	Type     string `json:"type"` // text | notification
}

type PaymentCreatedData struct {
	PaymentID  string  `json:"paymentId"`
	OrderID    string  `json:"orderId"`
	ProphetID  string  `json:"prophetId"`
	CourseID   string  `json:"courseId"`
	CustomerID string  `json:"customerId"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
}

type ChatMessageOutgoingData struct {
	MessageID   string `json:"messageId"`
	RoomID      string `json:"roomId"`
	SenderID    string `json:"senderId"`
	Content     string `json:"content"`
	Type        string `json:"type"` // text | notification
	CreatedAt   string `json:"createdAt"`
}
