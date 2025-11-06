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
	NotifyOrderCompleted     = "notify_order_completed"
)

// ---- DATA STRUCTURES ----

type OrderData struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
}

type OrderCompletedData struct {
	OrderID    string `json:"order_id"`
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
	OrderStatus string `json:"order_status"`
	ProphetID  string `json:"prophet_id"`
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

type PaymentPublishedData struct {
	PaymentID  string  `json:"paymentId"`
	OrderID    string  `json:"orderId"`
	ProphetID  string  `json:"prophetId"`
	CourseID   string  `json:"courseId"`
	CustomerID string  `json:"customerId"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
}

type ChatMessageOutgoingData struct {
	MessageID string `json:"messageId"`
	RoomID    string `json:"roomId"`
	SenderID  string `json:"senderId"`
	Content   string `json:"content"`
	Type      string `json:"type"` // text | notification
	CreatedAt string `json:"createdAt"`
}

type OrderCompletedData struct {
	OrderID     string `json:"orderId"`
	PaymentID   string `json:"paymentId"`
	OrderStatus string `json:"orderStatus"`
	CourseID    string `json:"courseId"`
	CourseName  string `json:"courseName"`
	RoomID      string `json:"roomId"`
}

type OrderPaymentBoundData struct {
	OrderID       string  `json:"orderId"`
	PaymentID     string  `json:"paymentId"`
	RoomID        string  `json:"roomId"`
	CustomerID    string  `json:"customerId"`
	OrderStatus   string  `json:"orderStatus"`
	CourseID      string  `json:"courseId"`
	CourseName    string  `json:"courseName"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"paymentStatus"`
}

type OrderPaidData struct {
	OrderID       string  `json:"orderId"`
	PaymentID     string  `json:"paymentId"`
	RoomID        string  `json:"roomId"`
	CustomerID    string  `json:"customerId"`
	CourseID      string  `json:"courseId"`
	OrderStatus   string  `json:"orderStatus"`
	CourseName    string  `json:"courseName"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"paymentStatus"`
}

type ChatNotificationOutgoingData[T any] struct {
	MessageID     string `json:"messageId"`
	RoomID        string `json:"roomId"`
	SenderID      string `json:"senderId"`
	Type          string `json:"type"` // text | notification
	CreatedAt     string `json:"createdAt"`
	MessageDetail *T     `json:"messageDetail"`
	Trigger       string `json:"trigger"`
}

type OrderCompletedNotificationData struct {
	OrderID     string `json:"orderId"`
	CourseID    string `json:"courseId"`
	OrderStatus string `json:"orderStatus"`
	CourseName  string `json:"courseName"`
}

type OrderPaymentBoundNotificationData struct {
	OrderID       string  `json:"orderId"`
	PaymentID     string  `json:"paymentId"`
	RoomID        string  `json:"roomId"`
	CustomerID    string  `json:"customerId"`
	CourseID      string  `json:"courseId"`
	OrderStatus   string  `json:"orderStatus"`
	CourseName    string  `json:"courseName"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"paymentStatus"`
}

type OrderPaidNotificationData struct {
	OrderID       string  `json:"orderId"`
	PaymentID     string  `json:"paymentId"`
	RoomID        string  `json:"roomId"`
	CustomerID    string  `json:"customerId"`
	CourseID      string  `json:"courseId"`
	OrderStatus   string  `json:"orderStatus"`
	CourseName    string  `json:"courseName"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"paymentStatus"`
}
