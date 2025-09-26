package messaging

import (
	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient manages a persistent RabbitMQ connection and channel.
type RabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client with a persistent connection and channel.
//
// Parameters:
//   - url: AMQP connection string (e.g., "amqp://guest:guest@localhost:5672/").
//
// Example:
//
//	client, err := NewRabbitMQClient("amqp://localhost:5672/")
func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQClient{conn: conn, channel: ch}, nil
}

// PublishToQueue publishes a message to a RabbitMQ queue.
//
// Parameters:
//   - queue: Queue name.
//   - message: The message body.
//
// Example:
//
//	err := client.PublishToQueue("jobs", "task1")
func (c *RabbitMQClient) PublishToQueue(queue, message string) error {
	_, err := c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return c.channel.Publish("", queue, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

// ConsumeFromQueue consumes messages from a RabbitMQ queue and handles each message using the provided handler.
//
// Parameters:
//   - queue: Queue name.
//   - handler: Function to handle each message body.
//
// Example:
//
//	err := client.ConsumeFromQueue("jobs", func(msg string) {
//	    fmt.Println("Received:", msg)
//	})
func (c *RabbitMQClient) ConsumeFromQueue(queue string, handler func(string)) error {
	_, err := c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		handler(string(msg.Body))
	}

	return nil
}

// Close closes the RabbitMQ connection and channel.
//
// Example:
//
//	err := client.Close()
func (c *RabbitMQClient) Close() error {
	if err := c.channel.Close(); err != nil {
		c.conn.Close()
		return err
	}
	return c.conn.Close()
}
