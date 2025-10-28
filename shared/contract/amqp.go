package contract

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

// Routing keys - using consistent event/command patterns
const (
	OrderCreatedEvent  = "order.created"
	PaymentSuccessEvent = "payment.complete"
	PaymentCreatedEvent = "payment.created"
	ChatMessageIncomingEvent = "chat.message.incoming"
	ChatMessageOutgoingEvent = "chat.message.outgoing"
)