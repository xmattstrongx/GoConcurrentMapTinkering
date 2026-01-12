package goconcurrentmaptinkering

import (
	"fmt"
	"testing"
	"time"

	cmap "github.com/xmattstrongx/go_concurrent_map"
)

func BenchmarkConcurrentMap_ReadWrite(b *testing.B) {
	m, err := cmap.New().
		WithDefaultExpiration(2 * time.Second).
		WithPurgeInterval(500 * time.Millisecond).
		Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}

	keys := make([]string, 256)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%len(keys)]
			m.Set(key, []byte("value"))
			m.Get(key)
			i++
		}
	})
}

func BenchmarkConcurrentMap_ReadHeavy(b *testing.B) {
	m, err := cmap.New().
		WithDefaultExpiration(0).
		WithPurgeInterval(1 * time.Second).
		Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}

	keys := make([]string, 1024)
	for i := range keys {
		key := fmt.Sprintf("key-%d", i)
		keys[i] = key
		m.Set(key, []byte("seed"))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(keys[i%len(keys)])
			m.Get(keys[(i+13)%len(keys)])
			m.Get(keys[(i+101)%len(keys)])
			i++
		}
	})
}

func BenchmarkConcurrentMap_WriteHeavy(b *testing.B) {
	m, err := cmap.New().
		WithDefaultExpiration(0).
		WithPurgeInterval(1 * time.Second).
		Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%len(keys)]
			m.Set(key, []byte("value"))
			i++
		}
	})
}
