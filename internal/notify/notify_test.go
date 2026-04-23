package notify_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/portwatch/internal/notify"
)

func newTestServer(t *testing.T, status int, received *notify.Payload) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if received != nil {
			if err := json.NewDecoder(r.Body).Decode(received); err != nil {
				http.Error(w, "bad body", http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(status)
	}))
}

func TestWebhookNotifier_Success(t *testing.T) {
	var got notify.Payload
	srv := newTestServer(t, http.StatusOK, &got)
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL, 0)
	p := notify.Payload{
		Event:   "port_opened",
		Port:    8080,
		Proto:   "tcp",
		Message: "new port detected",
	}

	if err := n.Notify(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Port != 8080 {
		t.Errorf("port: want 8080, got %d", got.Port)
	}
	if got.Event != "port_opened" {
		t.Errorf("event: want port_opened, got %s", got.Event)
	}
}

func TestWebhookNotifier_NonOKStatus(t *testing.T) {
	srv := newTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL, 0)
	err := n.Notify(context.Background(), notify.Payload{})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestWebhookNotifier_Timeout(t *testing.T) {
	blocking := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer blocking.Close()

	n := notify.NewWebhookNotifier(blocking.URL, 50*time.Millisecond)
	err := n.Notify(context.Background(), notify.Payload{})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestWebhookNotifier_InvalidURL(t *testing.T) {
	n := notify.NewWebhookNotifier("http://127.0.0.1:0", 100*time.Millisecond)
	err := n.Notify(context.Background(), notify.Payload{Port: 22})
	if err == nil {
		t.Fatal("expected connection error for invalid URL")
	}
}
