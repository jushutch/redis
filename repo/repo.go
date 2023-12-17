package repo

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// Repo manages the storing and retreiving of data
type Repo struct {
	logger      *slog.Logger
	data        *ThreadSafeMap[string]
	expirations *ThreadSafeMap[int64]
}

// New creates a new repo
func New(logger *slog.Logger) *Repo {
	return &Repo{
		logger:      logger.With("name", "redis.repo"),
		data:        NewThreadSafeMap[string](),
		expirations: NewThreadSafeMap[int64](),
	}
}

// Set writes a given value for a given key
func (r *Repo) Set(key, value string, expiration int64) error {
	r.logger.Info("set value", "key", key, "value", value, "expiration", expiration)
	r.data.Set(key, value)
	r.expirations.Set(key, expiration)
	return nil
}

// Get retrieves the value for a given key
func (r *Repo) Get(key string) (string, error) {
	r.logger.Info("get value", "key", key)
	if expirationUnix, err := r.expirations.Get(key); err != nil || isExpired(expirationUnix) {
		return "", fmt.Errorf("the key %q expired", key)
	}
	return r.data.Get(key)
}

// Delete removes the given key
func (r *Repo) Delete(key string) error {
	r.logger.Info("delete key", "key", key)
	if _, err := r.data.Get(key); err != nil {
		return fmt.Errorf("failed to get value for key: %w", err)
	}
	r.expirations.Set(key, -1)
	return nil
}

// Set writes a given value for a given key
func (r *Repo) Increment(key string) (int64, error) {
	r.logger.Info("increment", "key", key)
	valueStr, err := r.Get(key)
	if err != nil {
		err = r.Set(key, "1", 0)
		if err != nil {
			return 0, fmt.Errorf("failed to set new value: %w", err)
		}
		return 1, nil
	}
	value, err := strconv.ParseInt(valueStr, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer: %w", err)
	}
	value++
	r.data.Set(key, fmt.Sprintf("%d", value))
	return value, nil
}

func isExpired(unixMilli int64) bool {
	return unixMilli != 0 && (unixMilli == -1 || time.Now().After(time.UnixMilli(unixMilli)))
}
