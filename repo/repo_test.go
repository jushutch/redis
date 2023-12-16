package repo

import (
	"log/slog"
	"os"
	"testing"
)

func TestRepo(t *testing.T) {
	tests := []struct {
		name string
		test func(*Repo)
	}{
		{
			name: "set and get",
			test: func(r *Repo) {
				err := r.Set("key", "value")
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != "value" {
					t.Fatalf("got %s, wanted %s", value, "value")
				}
			},
		},
		{
			name: "set and get one key twice",
			test: func(r *Repo) {
				err := r.Set("key", "value")
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != "value" {
					t.Fatalf("got %s, wanted %s", value, "value")
				}
				err = r.Set("key", "new value")
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
				err := r.Set("key", "value")
				if err != nil {
					t.Fatalf("unexpected error setting value: %v", err)
				}
				value, err := r.Get("key")
				if err != nil {
					t.Fatalf("unexpected error getting value: %v", err)
				}
				if value != "value" {
					t.Fatalf("got %s, wanted %s", value, "value")
				}
				err = r.Set("second key", "second value")
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
