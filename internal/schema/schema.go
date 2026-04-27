// Package schema provides versioned serialisation helpers for portwatch
// persistent data files (baseline, history, audit log, etc.).
//
// Every file written to disk is wrapped in a thin envelope that records the
// schema version and the application version that wrote it.  Readers can
// inspect the version field before attempting to decode the payload, making
// forward- and backward-compatibility checks straightforward.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// Current is the schema version produced by this build.
const Current = 1

// ErrVersionMismatch is returned when a file was written with a schema version
// that this binary cannot safely decode.
var ErrVersionMismatch = errors.New("schema: unsupported version")

// Header carries the metadata that is prepended to every portwatch data file.
type Header struct {
	// Version is the schema version used when the file was written.
	Version int `json:"version"`

	// WrittenAt is the UTC timestamp at which the file was written.
	WrittenAt time.Time `json:"written_at"`

	// AppVersion is the portwatch release string (e.g. "v0.9.1"), if known.
	AppVersion string `json:"app_version,omitempty"`
}

// File is the top-level wrapper stored on disk.
type File[T any] struct {
	Header
	Payload T `json:"payload"`
}

// Write serialises payload into a versioned File and writes it to w.
// appVersion may be empty when the build version is unavailable.
func Write[T any](w io.Writer, payload T, appVersion string) error {
	f := File[T]{
		Header: Header{
			Version:    Current,
			WrittenAt:  time.Now().UTC(),
			AppVersion: appVersion,
		},
		Payload: payload,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(f); err != nil {
		return fmt.Errorf("schema: encode: %w", err)
	}
	return nil
}

// Read deserialises a versioned File from r into dst.
// It returns ErrVersionMismatch when the stored version is newer than Current.
func Read[T any](r io.Reader, dst *T) (Header, error) {
	var f File[T]
	if err := json.NewDecoder(r).Decode(&f); err != nil {
		return Header{}, fmt.Errorf("schema: decode: %w", err)
	}

	if f.Version > Current {
		return f.Header, fmt.Errorf("%w: file=%d current=%d",
			ErrVersionMismatch, f.Version, Current)
	}

	*dst = f.Payload
	return f.Header, nil
}

// CheckVersion decodes only the header from r and returns it without
// attempting to decode the payload.  Useful for pre-flight validation.
func CheckVersion(r io.Reader) (Header, error) {
	var h struct {
		Header
	}
	if err := json.NewDecoder(r).Decode(&h); err != nil {
		return Header{}, fmt.Errorf("schema: read header: %w", err)
	}
	return h.Header, nil
}
