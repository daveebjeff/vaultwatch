package notify

import (
	"testing"
	"time"
)

func TestMessage_Fields(t *testing.T) {
	now := time.Now()
	msg := Message{
		SecretPath: "secret/myapp/db",
		Status:     StatusExpiringSoon,
		Expiry:     now,
		Detail:     "expires in 2h",
	}

	if msg.SecretPath != "secret/myapp/db" {
		t.Errorf("unexpected SecretPath: %s", msg.SecretPath)
	}
	if msg.Status != StatusExpiringSoon {
		t.Errorf("unexpected Status: %s", msg.Status)
	}
	if !msg.Expiry.Equal(now) {
		t.Errorf("unexpected Expiry")
	}
	if msg.Detail != "expires in 2h" {
		t.Errorf("unexpected Detail: %s", msg.Detail)
	}
}

func TestStatus_Constants(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{StatusExpired, "EXPIRED"},
		{StatusExpiringSoon, "EXPIRING_SOON"},
		{StatusOK, "OK"},
	}
	for _, c := range cases {
		if string(c.s) != c.want {
			t.Errorf("Status %q: got %q", c.want, c.s)
		}
	}
}
