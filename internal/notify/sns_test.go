package notify

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type mockSNSPublisher struct {
	called bool
	input  *sns.PublishInput
	err    error
}

func (m *mockSNSPublisher) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.called = true
	m.input = params
	return &sns.PublishOutput{}, m.err
}

func TestNewSNSNotifier_EmptyARN(t *testing.T) {
	_, err := NewSNSNotifier("")
	if err == nil {
		t.Fatal("expected error for empty ARN")
	}
}

func TestSNSNotifier_Send_Success(t *testing.T) {
	mock := &mockSNSPublisher{}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123456789012:alerts"}

	msg := Message{
		SecretPath: "secret/db",
		Status:     StatusExpiringSoon,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
		Detail:     "expires in 10m",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected Publish to be called")
	}
	if mock.input == nil || *mock.input.TopicArn != n.topicARN {
		t.Error("unexpected topic ARN in publish input")
	}
}

func TestSNSNotifier_Send_Failure(t *testing.T) {
	mock := &mockSNSPublisher{err: errors.New("publish error")}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123456789012:alerts"}

	msg := Message{SecretPath: "secret/db", Status: StatusExpired, Detail: "expired"}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error from failed publish")
	}
}
