// Package envelope provides the Wrapper type, which decorates any alert
// payload with delivery metadata before it leaves the portwatch process.
//
// # Motivation
//
// When portwatch forwards events to external systems (webhooks, log
// aggregators, SIEM pipelines) each message needs enough context to be
// useful in isolation:
//
//   - Which host emitted the alert?
//   - Which daemon process (in case of restarts)?
//   - In what order did events arrive?
//   - When exactly did the alert occur?
//
// Rather than repeating this logic in every notifier, the pipeline wraps
// every outbound event once with [Wrapper.Wrap] and notifiers receive an
// [Envelope] that already carries the answers.
//
// # Usage
//
//	w := envelope.New()
//	env := w.Wrap(myEvent)      // Envelope{Seq:1, Hostname:"box", ...}
//	fmt.Println(env.String())   // [#1 box pid=12345 2024-06-01T12:00:00Z]
package envelope
