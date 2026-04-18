package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTeamsNotifier_EmptyURL(t *testing.T) {
	_, err := NewTeamsNotifier("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewTeamsNotifier_Valid(t *testing.T) {
	n, err := NewTeamsNotifier("https://outlook.office.com/webhook/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTeamsNotifier_Send_Success(t *testing.T) {
	var received bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewTeamsNotifier(ts.URL)
	msg := Message{
		Summary:    "secret expiring",
		SecretPath: "secret/db",
		Status:     StatusExpiringSoon,
		ExpiresAt:  time.Now().Add(24 * 60 * 60 * 1e9),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !received {
		t.Fatal("server did not receive request")
	}
}

func TestTeamsNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewTeamsNotifier(ts.URL)
	msg := Message{
		Summary:    "expired secret",
		SecretPath: "secret/api",
		Status:     StatusExpired,
		ExpiresAt:  time.Now().Add(-1),
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestThemeColor(t *testing.T) {
	tests := []struct {
		status Status
		color  string
	}{
		{StatusExpired, "FF0000"},
		{StatusExpiringSoon, "FFA500"},
		{StatusOK, "00FF00"},
	}
	for _, tc := range tests {
		if got := themeColor(tc.status); got != tc.color {
			t.Errorf("themeColor(%s) = %s, want %s", tc.status, got, tc.color)
		}
	}
}
