package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusPageHandler_EmptySnapshots(t *testing.T) {
	sn, _ := NewSnapshotNotifier(NewNoopNotifier())
	h := StatusPageHandler(sn)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []interface{}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d entries", len(out))
	}
}

func TestStatusPageHandler_WithSnapshots(t *testing.T) {
	sn, _ := NewSnapshotNotifier(NewNoopNotifier())
	_ = sn.Send(Message{Path: "secret/db", Status: StatusExpiringSoon})
	_ = sn.Send(Message{Path: "secret/api", Status: StatusOK})

	h := StatusPageHandler(sn)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e["path"] == nil || e["status"] == nil || e["received_at"] == nil {
			t.Errorf("missing fields in entry: %v", e)
		}
	}
}

func TestStatusPageHandler_ContentType(t *testing.T) {
	sn, _ := NewSnapshotNotifier(NewNoopNotifier())
	h := StatusPageHandler(sn)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
