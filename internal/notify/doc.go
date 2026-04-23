// Package notify implements webhook delivery for portwatch change events.
//
// A WebhookNotifier marshals a Payload to JSON and HTTP-POSTs it to a
// caller-supplied URL. It is intended to be composed with the alert pipeline:
//
//	n := notify.NewWebhookNotifier(cfg.Notify.WebhookURL, 5*time.Second)
//	p := notify.Payload{
//	    Event:     "port_opened",
//	    Port:      event.Port,
//	    Proto:     event.Proto,
//	    Timestamp: time.Now().UTC().Format(time.RFC3339),
//	    Message:   event.Message,
//	}
//	if err := n.Notify(ctx, p); err != nil {
//	    log.Printf("webhook delivery failed: %v", err)
//	}
//
// The notifier respects context cancellation and enforces a configurable
// per-request timeout so that a slow or unresponsive endpoint never blocks
// the main scan loop.
package notify
