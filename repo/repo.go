package repo

import (
	"fmt"
	"log/slog"
	"sync"
)

// Repo manages the storing and retreiving of data
type Repo struct {
	logger *slog.Logger
	data   map[string]string
	mu     sync.RWMutex
}

// New creates a new repo
func New(logger *slog.Logger) *Repo {
	return &Repo{
		logger: logger.With("name", "redis.repo"),
		data:   make(map[string]string),
		mu:     sync.RWMutex{},
	}
}

// Set writes a given value for a given key
func (r *Repo) Set(key, value string) error {
	r.mu.Lock()
	r.logger.Info("set value", "key", key, "value", value)
	r.data[key] = value
	r.mu.Unlock()
	return nil
}

// Get retrieves the value for a given key
func (r *Repo) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.logger.Info("get value", "key", key)
	value, ok := r.data[key]
	if !ok {
		return "", fmt.Errorf("the key %s does not exist", key)
	}
	return value, nil
}
