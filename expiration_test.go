package goconcurrentmaptinkering

import (
	"context"
	"testing"
	"time"

	cmap "github.com/xmattstrongx/go_concurrent_map"
)

func newTestMap(t *testing.T, defaultExp, purgeInterval time.Duration) *cmap.Concurrentmap {
	t.Helper()
	m, err := cmap.New().
		WithDefaultExpiration(defaultExp).
		WithPurgeInterval(purgeInterval).
		Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	return m
}

func TestExpiration_DefaultAndPerEntry(t *testing.T) {
	m := newTestMap(t, 80*time.Millisecond, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		m.PurgeExpiredEntries(ctx)
		close(done)
	}()

	m.Set("default", []byte("value"))
	m.SetEntry("short", cmap.Entry{
		KeyExpiration: 30 * time.Millisecond,
		Value:         []byte("short"),
	})
	t.Cleanup(func() {
		m.Delete("default")
		m.Delete("short")
	})

	deadline := time.Now().Add(400 * time.Millisecond)
	for {
		_, defaultOK := m.Get("default")
		_, shortOK := m.Get("short")
		if !shortOK && !defaultOK {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("entries did not expire before deadline (default=%v short=%v)", defaultOK, shortOK)
		}
		time.Sleep(5 * time.Millisecond)
	}

	<-done
}

func TestExpiration_NeverExpire(t *testing.T) {
	m := newTestMap(t, 40*time.Millisecond, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		m.PurgeExpiredEntries(ctx)
		close(done)
	}()

	m.SetEntry("forever", cmap.Entry{
		NeverExpire: true,
		Value:       []byte("keep"),
	})
	t.Cleanup(func() {
		m.Delete("forever")
	})

	time.Sleep(120 * time.Millisecond)
	if _, ok := m.Get("forever"); !ok {
		t.Fatalf("expected entry to remain when NeverExpire is set")
	}

	<-done
}
