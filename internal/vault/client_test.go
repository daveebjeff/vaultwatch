package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_InvalidAddress(t *testing.T) {
	// Should still succeed — Vault client defers connection errors.
	_, err := NewClient("http://127.0.0.1:1", "fake-token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}
}

func TestNewClient_EmptyToken(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:8200", "")
	if err != nil {
		t.Fatalf("did not expect error for empty token: %v", err)
	}
}

func TestIsAuthenticated_Failure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "bad-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.IsAuthenticated(); err == nil {
		t.Error("expected authentication error, got nil")
	}
}

func TestReadSecret_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ReadSecret("secret/missing")
	if err == nil {
		t.Error("expected error for missing secret")
	}
}
