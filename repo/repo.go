package repo

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Repo manages the storing and retreiving of data
type Repo struct {
	logger      *slog.Logger
	data        map[string]string
	expirations map[string]int64
	dmu         sync.RWMutex
	emu         sync.RWMutex
}

// New creates a new repo
func New(logger *slog.Logger) *Repo {
	return &Repo{
		logger:      logger.With("name", "redis.repo"),
		data:        make(map[string]string),
		expirations: make(map[string]int64),
	}
}

// Set writes a given value for a given key
func (r *Repo) Set(key, value string, expiration int64) error {
	// Save cycles if there are no changes
	if !r.isNewValue(key, value) && !r.isNewExpiration(key, expiration) {
		return nil
	}

	r.dmu.Lock()
	r.logger.Info("set value", "key", key, "value", value, "expiration", expiration)
	r.data[key] = value
	r.dmu.Unlock()
	r.emu.Lock()
	r.expirations[key] = expiration
	r.emu.Unlock()
	return nil
}

// Get retrieves the value for a given key
func (r *Repo) Get(key string) (string, error) {
	r.emu.RLock()
	if expirationUnix, ok := r.expirations[key]; ok && isExpired(expirationUnix) {
		r.emu.RUnlock()
		return "", fmt.Errorf("the key %q expired", key)
	}
	r.emu.RUnlock()

	r.dmu.RLock()
	defer r.dmu.RUnlock()
	r.logger.Info("get value", "key", key)
	value, ok := r.data[key]
	if !ok {
		return "", fmt.Errorf("the key %q does not exist", key)
	}
	return value, nil
}

func (r *Repo) isNewExpiration(key string, expiration int64) bool {
	r.emu.RLock()
	val, ok := r.expirations[key]
	new := ok && val == expiration
	r.emu.RUnlock()
	return new
}

func (r *Repo) isNewValue(key string, value string) bool {
	r.dmu.RLock()
	val, ok := r.data[key]
	new := ok && val == value
	r.dmu.RUnlock()
	return new
}

func isExpired(unixMilli int64) bool {
	return unixMilli != 0 && time.Now().After(time.UnixMilli(unixMilli))
}
