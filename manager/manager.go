package manager

import (
	"fmt"

	"github.com/jushutch/redis/repo"
	"github.com/jushutch/redis/serializer"
)

type Repo interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

type Manager struct {
	repo Repo
}

func NewManager() *Manager {
	return &Manager{
		repo: repo.NewRepo(),
	}
}

func (m *Manager) Ping() serializer.RESPType {
	return serializer.SimpleString("PONG")
}

func (m *Manager) Echo(message serializer.BulkString) serializer.RESPType {
	return message
}

func (m *Manager) Set(key, value serializer.BulkString) serializer.RESPType {
	err := m.repo.Set(key.Value, value.Value)
	if err != nil {
		return nil
	}
	return serializer.SimpleString("OK")
}

func (m *Manager) Get(key serializer.BulkString) serializer.RESPType {
	value, err := m.repo.Get(key.Value)
	if err != nil {
		fmt.Println(err.Error())
		return serializer.BulkString{Length: -1}
	}
	return serializer.BulkString{Length: int64(len(value)), Value: value}
}
