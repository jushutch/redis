package repo

import (
	"fmt"
	"sync"
)

type ThreadSafeMap[T any] struct {
	data map[string]T
	mu   sync.RWMutex
}

func NewThreadSafeMap[T any]() *ThreadSafeMap[T] {
	return &ThreadSafeMap[T]{
		data: make(map[string]T),
	}
}

func (m *ThreadSafeMap[T]) Set(key string, value T) {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
}

func (m *ThreadSafeMap[T]) Get(key string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	if !ok {
		return value, fmt.Errorf("the key %q does not exist", key)
	}
	return value, nil
}
