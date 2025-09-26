package messaging

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// PublishMessage publishes a message to the given Redis channel.
func PublishMessage(ctx context.Context, client *redis.Client, channel, message string) error {
	return client.Publish(ctx, channel, message).Err()
}

// SubscribeToChannel subscribes to a Redis channel and handles incoming messages with the provided handler.
func SubscribeToChannel(ctx context.Context, client *redis.Client, channel string, handler func(string)) error {
	sub := client.Subscribe(ctx, channel)
	ch := sub.Channel()

	for {
		select {
		case msg := <-ch:
			handler(msg.Payload)
		case <-ctx.Done():
			return sub.Close()
		}
	}
}
