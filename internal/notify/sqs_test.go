package notify

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type mockSQSSender struct {
	called bool
	input  *sqs.SendMessageInput
	err    error
}

func (m *mockSQSSender) SendMessage(_ context.Context, params *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.called = true
	m.input = params
	return &sqs.SendMessageOutput{}, m.err
}

func TestNewSQSNotifier_EmptyURL(t *testing.T) {
	_, err := NewSQSNotifier("")
	if err == nil {
		t.Fatal("expected error for empty queue URL")
	}
}

func TestSQSNotifier_Send_Success(t *testing.T) {
	mock := &mockSQSSender{}
	n := &SQSNotifier{client: mock, queueURL: "https://sqs.us-east-1.amazonaws.com/123/alerts"}

	msg := Message{
		SecretPath: "secret/api",
		Status:     StatusExpired,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
		Detail:     "already expired",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected SendMessage to be called")
	}
	if mock.input == nil || !strings.Contains(*mock.input.MessageBody, "secret/api") {
		t.Error("message body missing secret path")
	}
}

func TestSQSNotifier_Send_Failure(t *testing.T) {
	mock := &mockSQSSender{err: errors.New("queue error")}
	n := &SQSNotifier{client: mock, queueURL: "https://sqs.us-east-1.amazonaws.com/123/alerts"}

	msg := Message{SecretPath: "secret/api", Status: StatusExpired, Detail: "expired"}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error from failed send")
	}
}
