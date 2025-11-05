package contract

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

// Routing keys - using consistent event/command patterns
const (
	OrderCreatedEvent  = "order.created"
	OrderCompletedEvent = "order.completed"
	PaymentSuccessEvent = "payment.completed"
	PaymentCreatedEvent = "payment.created"
	PaymentSettledEvent = "payment.settled"
	ChatMessageIncomingEvent = "chat.message.incoming"
	ChatMessageOutgoingEvent = "chat.message.outgoing"
	OrderPaymentBoundEvent = "order.payment.bound"
	OrderPaidEvent = "order.paid"
)