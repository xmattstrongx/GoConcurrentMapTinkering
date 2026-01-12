package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	cmap "github.com/xmattstrongx/go_concurrent_map"
)

func main() {
	m, err := cmap.New().
		WithDefaultExpiration(2 * time.Second).
		WithPurgeInterval(250 * time.Millisecond).
		Build()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go m.PurgeExpiredEntries(ctx)
	defer cancel()

	m.Set("greeting", []byte("hello"))
	if v, ok := m.Get("greeting"); ok {
		fmt.Printf("greeting=%s\n", v)
	}

	m.SetEntry("short", cmap.Entry{
		KeyExpiration: 500 * time.Millisecond,
		Value:         []byte("bye"),
	})

	m.SetEntry("forever", cmap.Entry{
		NeverExpire: true,
		Value:       []byte("always"),
	})

	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}
	var wgWriters sync.WaitGroup
	var wgReaders sync.WaitGroup
	done := make(chan struct{})

	for i := 0; i < 2; i++ {
		wgWriters.Add(1)
		go func(id int) {
			defer wgWriters.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
			for j := 0; j < 60; j++ {
				key := keys[rng.Intn(len(keys))]
				exp := time.Duration(300+rng.Intn(700)) * time.Millisecond
				m.SetEntry(key, cmap.Entry{
					KeyExpiration: exp,
					Value:         []byte(fmt.Sprintf("payload-%d-%d", id, j)),
				})
				time.Sleep(30 * time.Millisecond)
			}
		}(i)
	}

	for i := 0; i < 4; i++ {
		wgReaders.Add(1)
		go func(id int) {
			defer wgReaders.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id+10)))
			for {
				select {
				case <-done:
					return
				default:
					key := keys[rng.Intn(len(keys))]
					if v, ok := m.Get(key); ok {
						_ = v
					}
					time.Sleep(15 * time.Millisecond)
				}
			}
		}(i)
	}

	wgWriters.Wait()
	close(done)
	wgReaders.Wait()

	time.Sleep(1 * time.Second)
	if _, ok := m.Get("short"); !ok {
		fmt.Println("short entry expired as expected")
	}
	if v, ok := m.Get("forever"); ok {
		fmt.Printf("forever=%s\n", v)
	}
}
