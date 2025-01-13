package pingo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestStartPinger(t *testing.T) {
	t.Run("starts with default options", func(t *testing.T) {
		var pingCount atomic.Int64
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pingCount.Add(1)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cancelPinger, err := StartPinger(ctx, WithURL(ts.URL))
		if err != nil {
			t.Fatalf("StartPinger failed: %v", err)
		}
		defer cancelPinger()

		time.Sleep(100 * time.Millisecond)

		if pingCount.Load() == 0 {
			t.Error("expected at least one ping")
		}
	})

	t.Run("respects custom frequency", func(t *testing.T) {
		var pingCount atomic.Int64
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pingCount.Add(1)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		customFreq := 100 * time.Millisecond
		cancelPinger, err := StartPinger(ctx, WithFrequency(customFreq), WithURL(ts.URL))
		if err != nil {
			t.Fatalf("StartPinger failed: %v", err)
		}
		defer cancelPinger()

		time.Sleep(250 * time.Millisecond)

		if pingCount.Load() < 2 {
			t.Errorf("expected at least 2 pings, got %d", pingCount.Load())
		}
	})

	t.Run("uses custom http client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 1 * time.Second,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cancelPinger, err := StartPinger(ctx, WithHTTPClient(customClient))
		if err != nil {
			t.Fatalf("StartPinger failed: %v", err)
		}
		defer cancelPinger()
	})

	t.Run("cancellation stops pinger", func(t *testing.T) {
		var pingCount atomic.Int64
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pingCount.Add(1)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		ctx, cancel := context.WithCancel(context.Background())

		cancelPinger, err := StartPinger(ctx, WithFrequency(50*time.Millisecond), WithURL(ts.URL))
		if err != nil {
			t.Fatalf("StartPinger failed: %v", err)
		}
		defer cancelPinger()

		time.Sleep(200 * time.Millisecond)
		initialCount := pingCount.Load()

		cancel()

		time.Sleep(200 * time.Millisecond)

		if pingCount.Load() > initialCount+1 {
			t.Errorf("pinger continued after cancellation")
		}
	})
}
