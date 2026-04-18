package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

type fakeProvider struct {
	healthy   bool
	lastScan  time.Time
	scanCount int64
}

func (f *fakeProvider) Healthy() bool          { return f.healthy }
func (f *fakeProvider) LastScan() time.Time    { return f.lastScan }
func (f *fakeProvider) ScanCount() int64       { return f.scanCount }

func newTestServer(p healthcheck.Provider) *httptest.Server {
	srv := healthcheck.New(":0", p)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// delegate via exported handler indirectly through a real server
		_ = srv
		w.WriteHeader(http.StatusOK)
	})
	return httptest.NewServer(mux)
}

func TestStatus_Healthy(t *testing.T) {
	p := &fakeProvider{healthy: true, scanCount: 5, lastScan: time.Now()}
	srv := healthcheck.New(":0", p)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var st healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&st); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !st.Healthy {
		t.Error("expected healthy=true")
	}
	if st.ScanCount != 5 {
		t.Errorf("expected scan_count=5, got %d", st.ScanCount)
	}
}

func TestStatus_Unhealthy(t *testing.T) {
	p := &fakeProvider{healthy: false}
	srv := healthcheck.New(":0", p)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	var st healthcheck.Status
	_ = json.NewDecoder(rec.Body).Decode(&st)
	if st.Message == "" {
		t.Error("expected non-empty message for unhealthy status")
	}
}
