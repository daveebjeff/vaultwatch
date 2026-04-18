package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type sqsSender interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// SQSNotifier sends alert messages to an AWS SQS queue.
type SQSNotifier struct {
	client   sqsSender
	queueURL string
}

// NewSQSNotifier creates an SQSNotifier using the default AWS config.
func NewSQSNotifier(queueURL string) (*SQSNotifier, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("sqs: queue URL must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sqs: failed to load AWS config: %w", err)
	}
	return &SQSNotifier{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

// Send enqueues a JSON-encoded notification message to the configured SQS queue.
func (n *SQSNotifier) Send(msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("sqs: failed to marshal message: %w", err)
	}
	_, err = n.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(n.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("sqs: send failed: %w", err)
	}
	return nil
}
