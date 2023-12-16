package repo

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

const (
	mockKey   = "key"
	mockValue = "value"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		test func(*Repo)
	}{
		{
			name: "set and get",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, 0)
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
			},
		},
		{
			name: "set and get one key twice",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, 0)
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
				err = r.Set("key", "new value", 0)
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err = r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != "new value" {
					t.Fatalf("got %q, wanted %q", value, "new value")
				}
			},
		},
		{
			name: "get unset key",
			test: func(r *Repo) {
				value, err := r.Get("key")
				if err == nil {
					t.Fatalf("expected error getting value")
				}
				if value != "" {
					t.Fatalf("expected value to be empty")
				}
			},
		},
		{
			name: "set and get two keys",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, 0)
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
				err = r.Set("second key", "second value", 0)
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err = r.Get("second key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != "second value" {
					t.Fatalf("got %q, wanted %q", value, "second value")
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := slog.New(handler)
			repo := New(logger)
			test.test(repo)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		test func(*Repo)
	}{
		{
			name: "get unset key",
			test: func(r *Repo) {
				value, err := r.Get("key")
				if err == nil {
					t.Fatalf("expected error getting value")
				}
				if value != "" {
					t.Fatalf("expected value to be empty")
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := slog.New(handler)
			repo := New(logger)
			test.test(repo)
		})
	}
}

func TestExpiration(t *testing.T) {
	tests := []struct {
		name string
		test func(*Repo)
	}{
		{
			name: "unexpired",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, time.Now().Add(10*time.Second).UnixMilli())
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
			},
		},
		{
			name: "expired",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, time.Now().Add(time.Second).UnixMilli())
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				time.Sleep(time.Second)
				value, err := r.Get("key")
				if err == nil {
					t.Fatalf("expected error for expired key")
				}
				if value != "" {
					t.Fatalf("expected empty value for expired key")
				}
			},
		},
		{
			name: "before and after expiration",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, time.Now().Add(time.Second).UnixMilli())
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
				time.Sleep(time.Second)
				value, err = r.Get("key")
				if err == nil {
					t.Fatalf("expected error for expired key")
				}
				if value != "" {
					t.Fatalf("expected empty value for expired key")
				}
			},
		},
		{
			name: "set expiration on existing key",
			test: func(r *Repo) {
				err := r.Set("key", mockValue, time.Now().Add(time.Second).UnixMilli())
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
				time.Sleep(time.Second)
				value, err = r.Get("key")
				if err == nil {
					t.Fatalf("expected error for expired key")
				}
				if value != "" {
					t.Fatalf("expected empty value for expired key")
				}
				err = r.Set("key", mockValue, 0)
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err = r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != mockValue {
					t.Fatalf("got %s, wanted %s", value, mockValue)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := slog.New(handler)
			repo := New(logger)
			test.test(repo)
		})
	}
}
