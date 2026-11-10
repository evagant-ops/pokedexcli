package pokecache

import (
	"context"
	"reflect"
	"testing"
	"time"
)

const testInterval = time.Hour

// newTestCache returns a cache whose reaper is stopped when the test ends.
func newTestCache(t *testing.T, interval time.Duration) *Cache {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return NewCache(ctx, interval)
}

func TestNewCache(t *testing.T) {
	c := newTestCache(t, testInterval)
	if c == nil {
		t.Fatal("NewCache() returned nil")
	}
	if c.Data == nil {
		t.Fatal("NewCache() Data map is nil")
	}
	if len(c.Data) != 0 {
		t.Errorf("NewCache() Data len = %d, want 0", len(c.Data))
	}
}

func TestCache_Add(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  []byte
	}{
		{"simple value", "key1", []byte("value1")},
		{"empty value", "empty", []byte{}},
		{"nil value", "nil", nil},
		{"empty key", "", []byte("value")},
		{"url key", "https://pokeapi.co/api/v2/location-area", []byte(`{"count":1}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCache(t, testInterval)
			c.Add(tt.key, tt.val)
			e, ok := c.Data[tt.key]
			got := e.val
			if !ok {
				t.Fatalf("key %q not stored", tt.key)
			}
			if !reflect.DeepEqual(got, tt.val) {
				t.Errorf("stored = %v, want %v", got, tt.val)
			}
		})
	}
}

func TestCache_Add_Overwrite(t *testing.T) {
	c := newTestCache(t, testInterval)
	c.Add("key", []byte("old"))
	c.Add("key", []byte("new"))

	got, ok := c.Get("key")
	if !ok {
		t.Fatal("key missing after overwrite")
	}
	if string(got) != "new" {
		t.Errorf("got %q, want %q", got, "new")
	}
	if len(c.Data) != 1 {
		t.Errorf("Data len = %d, want 1", len(c.Data))
	}
}

func TestCache_Get(t *testing.T) {
	c := newTestCache(t, testInterval)
	c.Add("a", []byte("alpha"))
	c.Add("b", []byte("beta"))
	c.Add("empty", []byte{})

	tests := []struct {
		name   string
		key    string
		want   []byte
		wantOk bool
	}{
		{"existing key a", "a", []byte("alpha"), true},
		{"existing key b", "b", []byte("beta"), true},
		{"existing empty value", "empty", []byte{}, true},
		{"missing key", "missing", nil, false},
		{"missing empty key", "", nil, false},
		{"case sensitive key", "A", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := c.Get(tt.key)
			if ok != tt.wantOk {
				t.Errorf("Get(%q) ok = %v, want %v", tt.key, ok, tt.wantOk)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestCache_Get_EmptyCache(t *testing.T) {
	c := newTestCache(t, testInterval)
	got, ok := c.Get("anything")
	if ok {
		t.Error("Get on empty cache returned ok = true")
	}
	if got != nil {
		t.Errorf("Get on empty cache = %v, want nil", got)
	}
}

func TestCache_evictExpired(t *testing.T) {
	interval := time.Minute
	now := time.Now()

	tests := []struct {
		name     string
		entries  map[string]time.Time // key -> createdAt
		wantKeys []string
	}{
		{"empty cache", map[string]time.Time{}, nil},
		{"fresh entry kept", map[string]time.Time{"a": now}, []string{"a"}},
		{"old entry removed", map[string]time.Time{"a": now.Add(-2 * interval)}, nil},
		{
			"mixed: only old removed",
			map[string]time.Time{
				"old1":  now.Add(-time.Hour),
				"old2":  now.Add(-2 * interval),
				"fresh": now.Add(-interval / 2),
			},
			[]string{"fresh"},
		},
		{"all old removed", map[string]time.Time{"a": now.Add(-time.Hour), "b": now.Add(-time.Hour)}, nil},
		{"all fresh kept", map[string]time.Time{"a": now, "b": now}, []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// literal: no reaper goroutine, so the test controls eviction
			c := &Cache{interval: interval, Data: make(map[string]cacheEntry)}
			for k, created := range tt.entries {
				c.Data[k] = cacheEntry{val: []byte(k), createdAt: created}
			}

			c.evictExpired()

			if len(c.Data) != len(tt.wantKeys) {
				t.Errorf("len = %d, want %d (data=%v)", len(c.Data), len(tt.wantKeys), c.Data)
			}
			for _, k := range tt.wantKeys {
				if _, ok := c.Get(k); !ok {
					t.Errorf("key %q should remain", k)
				}
			}
		})
	}
}

func TestCache_reapLoop_evictsAfterInterval(t *testing.T) {
	c := newTestCache(t, 20*time.Millisecond)
	c.Add("key", []byte("val"))

	if _, ok := c.Get("key"); !ok {
		t.Fatal("key should exist right after Add")
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, ok := c.Get("key"); !ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Error("key not evicted by reapLoop")
}

func TestCache_reapLoop_stopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache{interval: 10 * time.Millisecond, Data: make(map[string]cacheEntry)}

	done := make(chan struct{})
	go func() {
		c.reapLoop(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("reapLoop did not exit after context cancel")
	}
}
