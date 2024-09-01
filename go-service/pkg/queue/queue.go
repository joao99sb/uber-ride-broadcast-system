package queue

import (
	"fmt"
	"go-service/pkg/types"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type Queue struct {
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Connect() {
	user := os.Getenv("RABBITMQ_DEFAULT_USER")
	pass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	host := os.Getenv("RABBITMQ_DEFAULT_HOST")
	port := os.Getenv("RABBITMQ_DEFAULT_PORT")
	vhost := os.Getenv("RABBITMQ_DEFAULT_VHOST")

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, pass, host, port, vhost)

	conn, err := amqp.Dial(dsn)
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	q.conn = conn
	q.channel = ch
}

func (q *Queue) GetChannel() *amqp.Channel {
	return q.channel
}
func (q *Queue) GetConnection() *amqp.Connection {
	return q.conn
}
func (q *Queue) Close() {
	q.conn.Close()
	q.channel.Close()
}

func (q *Queue) GetQueue(queueName string) amqp.Queue {
	queue, err := q.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")
	return queue
}
func (q *Queue) StartConsuming(que amqp.Queue) <-chan []byte {

	out := make(chan []byte)
	msg, err := q.channel.Consume(
		que.Name,
		"go-worker",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msg {

			out <- []byte(d.Body)
		}

	}()
	return out
}

func (q *Queue) BuildDestinationQueue(queueName, exchangeName, routingKey string) *types.QueueDestination {
	return &types.QueueDestination{
		QueueName:    queueName,
		ExchangeName: exchangeName,
		RoutingKey:   routingKey,
	}
}
func (q *Queue) BuildDefaultDestinationQueue() *types.QueueDestination {
	return &types.QueueDestination{
		QueueName:    os.Getenv("RABBITMQ_CONSUMER_QUEUE"),
		ExchangeName: os.Getenv("RABBITMQ_DESTINATION"),
		RoutingKey:   os.Getenv("RABBITMQ_DESTINATION_ROUTING_KEY"),
	}
}

func (q *Queue) Notify(message []byte, dest *types.QueueDestination) {
	err := q.channel.Publish(
		dest.ExchangeName,
		dest.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	FailOnError(err, "Failed to publish a message")

}

func (q *Queue) CreateQueueAndBind(dest *types.QueueDestination) error {

	// Bind the queue to the exchange
	err := q.channel.QueueBind(
		dest.QueueName,    // queue name
		dest.RoutingKey,   // routing key
		dest.ExchangeName, // exchange name
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	log.Printf("Queue '%s' created and bound to exchange '%s' with routing key '%s'", dest.QueueName, dest.ExchangeName, dest.RoutingKey)
	return nil
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
