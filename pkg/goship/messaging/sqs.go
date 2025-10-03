package messaging

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// NewSQSClient creates a new AWS SQS client using provided credentials and region.
//
// Parameters:
//   - accessKey: AWS access key ID (e.g., "AKIA...").
//   - secretKey: AWS secret access key (e.g., "abcd1234...").
//   - region: AWS region (e.g., "us-west-2").
//
// Example:
//
//	client, err := NewSQSClient("AKIA...", "abcd...", "us-west-2")
func NewSQSClient(accessKey, secretKey, region string) (*sqs.Client, error) {
	cfg := aws.Config{
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		Region:      region,
	}
	return sqs.NewFromConfig(cfg), nil
}

// SendMessageToQueue sends a message to an SQS queue.
func SendMessageToQueue(ctx context.Context, client *sqs.Client, queueURL, message string) error {
	_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(message),
	})
	return err
}

// ReceiveMessagesFromQueue continuously receives and processes messages from an SQS queue.
func ReceiveMessagesFromQueue(ctx context.Context, client *sqs.Client, queueURL string, handler func(string) error) error {
	for {
		resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			return err
		}

		for _, msg := range resp.Messages {
			if err := handler(*msg.Body); err == nil {
				_, _ = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})
			}
		}
	}
}
