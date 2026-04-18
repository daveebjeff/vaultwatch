package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// snsPublisher abstracts the SNS API for testing.
type snsPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotifier sends alert messages to an AWS SNS topic.
type SNSNotifier struct {
	client   snsPublisher
	topicARN string
}

// NewSNSNotifier creates an SNSNotifier using the default AWS config.
func NewSNSNotifier(topicARN string) (*SNSNotifier, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: failed to load AWS config: %w", err)
	}
	return &SNSNotifier{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// Send publishes a notification message to the configured SNS topic.
func (n *SNSNotifier) Send(msg Message) error {
	body := fmt.Sprintf("[%s] %s — %s", msg.Status, msg.SecretPath, msg.Detail)
	_, err := n.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Message:  aws.String(body),
		Subject:  aws.String(fmt.Sprintf("VaultWatch: %s", msg.Status)),
	})
	if err != nil {
		return fmt.Errorf("sns: publish failed: %w", err)
	}
	return nil
}
