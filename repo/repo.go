package repo

import (
	"fmt"
	"log/slog"
	"sync"
)

type Repo struct {
	logger *slog.Logger
	data   map[string]string
	mu     sync.RWMutex
}

func NewRepo(logger *slog.Logger) *Repo {
	return &Repo{
		logger: logger.With("name", "redis.repo"),
		data:   make(map[string]string),
	}
}

func (r *Repo) Set(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("set value", "key", key, "value", value)
	r.data[key] = value
	return nil
}

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
