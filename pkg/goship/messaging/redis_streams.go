package messaging

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// AddToStream adds a message to a Redis Stream.
func AddToStream(ctx context.Context, client *redis.Client, stream string, values map[string]interface{}) error {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Err()
}

// ConsumeStream reads from a Redis Stream and calls a handler for each message.
func ConsumeStream(ctx context.Context, client *redis.Client, stream, group, consumer string, handler func(redis.XMessage)) error {
	// Create group if not exists
	_ = client.XGroupCreateMkStream(ctx, stream, group, "0").Err()

	for {
		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Block:    0,
			Count:    1,
		}).Result()
		if err != nil {
			return err
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				handler(msg)
				client.XAck(ctx, stream.Stream, group, msg.ID)
			}
		}
	}
}
