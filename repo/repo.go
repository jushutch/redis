package repo

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/jushutch/redis/logging"
)

// Repo manages the storing and retreiving of data
type Repo struct {
	logger      *slog.Logger
	data        *ThreadSafeMap[interface{}]
	expirations *ThreadSafeMap[int64]
}

// New creates a new repo
func New(logger *slog.Logger) *Repo {
	return &Repo{
		logger:      logger.With("name", "redis.repo"),
		data:        NewThreadSafeMap[interface{}](),
		expirations: NewThreadSafeMap[int64](),
	}
}

// Set writes a given value for a given key
func (r *Repo) Set(ctx context.Context, key, value string, expiration int64) error {
	r.logger.With(logging.FieldsFromContext(ctx)...).Info("set value", "key", key, "value", value, "expiration", expiration)
	r.data.Set(key, value)
	r.expirations.Set(key, expiration)
	return nil
}

// Get retrieves the value for a given key
func (r *Repo) Get(ctx context.Context, key string) (string, error) {
	r.logger.With(logging.FieldsFromContext(ctx)...).Info("get value", "key", key)
	if expirationUnix, err := r.expirations.Get(key); err != nil || isExpired(expirationUnix) {
		return "", fmt.Errorf("the key %q expired", key)
	}
	rawValue, err := r.data.Get(key)
	if err != nil {
		return "", fmt.Errorf("failed to get value for key %q: %w", key, err)
	}
	if value, ok := rawValue.(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("value (%v) for requested key (%q) was not a string", rawValue, key)
}

// Delete removes the given key
func (r *Repo) Delete(ctx context.Context, key string) error {
	r.logger.With(logging.FieldsFromContext(ctx)...).Info("delete key", "key", key)
	if _, err := r.data.Get(key); err != nil {
		return fmt.Errorf("failed to get value for key: %w", err)
	}
	r.expirations.Set(key, -1)
	return nil
}

// Add attempts to perform addition on the stored value
func (r *Repo) Add(ctx context.Context, key string, toAdd int64) (int64, error) {
	r.logger.With(logging.FieldsFromContext(ctx)...).Info("add", "key", key, "to_add", toAdd)
	valueStr, err := r.Get(ctx, key)
	if err != nil {
		err = r.Set(ctx, key, strconv.FormatInt(toAdd, 10), 0)
		if err != nil {
			return 0, fmt.Errorf("failed to set new value: %w", err)
		}
		return toAdd, nil
	}
	value, err := strconv.ParseInt(valueStr, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer: %w", err)
	}
	if math.MaxInt64-value > toAdd {
		return 0, fmt.Errorf("new value would be out of range")
	}
	value += toAdd
	r.data.Set(key, strconv.FormatInt(value, 10))
	return value, nil
}

func isExpired(unixMilli int64) bool {
	return unixMilli != 0 && (unixMilli == -1 || time.Now().After(time.UnixMilli(unixMilli)))
}
