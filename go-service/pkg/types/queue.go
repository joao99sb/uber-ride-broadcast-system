package types

type IQueue interface {
	Connect()
	Notify(message []byte, dest *QueueDestination)
	BuildDefaultDestinationQueue() *QueueDestination
}

type QueueDestination struct {
	QueueName    string
	ExchangeName string
	RoutingKey   string
}
