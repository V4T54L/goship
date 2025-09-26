package messaging

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer is a reusable Kafka writer.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer for a given broker and topic.
//
// Parameters:
//   - brokerAddr: Kafka broker address (e.g., "localhost:9092").
//   - topic: Kafka topic name.
//
// Example:
//
//	producer := NewKafkaProducer("localhost:9092", "events")
func NewKafkaProducer(brokerAddr, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}
}

// Publish sends a message to the configured Kafka topic.
//
// Parameters:
//   - message: The message payload.
//
// Example:
//
//	err := producer.Publish("hello world")
func (p *KafkaProducer) Publish(message string) error {
	return p.writer.WriteMessages(context.Background(),
		kafka.Message{Value: []byte(message)},
	)
}

// Close shuts down the Kafka producer.
//
// Example:
//
//	err := producer.Close()
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// KafkaConsumer consumes messages from a Kafka topic using a reusable reader.
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer.
//
// Parameters:
//   - brokerAddr: Kafka broker address (e.g., "localhost:9092").
//   - topic: Kafka topic name.
//   - groupID: Consumer group ID.
//
// Example:
//
//	consumer := NewKafkaConsumer("localhost:9092", "events", "group-1")
func NewKafkaConsumer(brokerAddr, topic, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &KafkaConsumer{reader: reader}
}

// Consume starts consuming messages from the Kafka topic and calls the handler for each message.
//
// Parameters:
//   - handler: Function to handle each message body.
//
// Example:
//
//	err := consumer.Consume(func(msg string) {
//	    fmt.Println("Received:", msg)
//	})
func (c *KafkaConsumer) Consume(handler func(string)) error {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		handler(string(m.Value))
	}
}

// Close closes the Kafka reader.
//
// Example:
//
//	err := consumer.Close()
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
