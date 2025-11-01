package inbound_port

type MessageConsumer interface {
	StartListening() error
}
