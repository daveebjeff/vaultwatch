package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPGetNotifier_EmptyURL(t *testing.T) {
	_, err := NewHTTPGetNotifier("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewHTTPGetNotifier_Valid(t *testing.T) {
	n, err := NewHTTPGetNotifier("http://example.com/notify")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestHTTPGetNotifier_Send_Success(t *testing.T) {
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewHTTPGetNotifier(ts.URL)
	msg := Message{
		Status:     StatusExpired,
		SecretPath: "secret/db",
		ExpiresAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery == "" {
		t.Error("expected query parameters to be sent")
	}
}

func TestHTTPGetNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewHTTPGetNotifier(ts.URL)
	err := n.Send(Message{Status: StatusExpiringSoon, SecretPath: "secret/api"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestHTTPGetNotifier_Send_BadURL(t *testing.T) {
	n := &HTTPGetNotifier{baseURL: "://bad url", client: &http.Client{}}
	err := n.Send(Message{Status: StatusExpired, SecretPath: "x"})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestHTTPGetNotifier_Send_QueryContainsSecretPath(t *testing.T) {
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewHTTPGetNotifier(ts.URL)
	msg := Message{
		Status:     StatusExpiringSoon,
		SecretPath: "secret/myapp",
		ExpiresAt:  time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery == "" {
		t.Fatal("expected query parameters")
	}
	if !contains(gotQuery, "secret") {
		t.Errorf("expected secret path in query, got: %s", gotQuery)
	}
}

// contains is a helper to check substring presence in query strings.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}
