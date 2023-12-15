package manager

import (
	"log/slog"

	"github.com/jushutch/redis/repo"
	"github.com/jushutch/redis/serializer"
)

// Define allowed commands
const (
	PING = "PING"
	ECHO = "ECHO"
	SET  = "SET"
	GET  = "GET"
)

type Repo interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

type Manager struct {
	repo   Repo
	logger *slog.Logger
}

func NewManager(logger *slog.Logger) *Manager {
	return &Manager{
		logger: logger.With("name", "redis.manager"),
		repo:   repo.NewRepo(logger),
	}
}

func (m *Manager) HandleCommand(command serializer.Array) serializer.RESPType {
	if command.Length <= 0 {
		return nil
	}

	name, ok := command.Elements[0].(serializer.BulkString)
	if !ok {
		return nil
	}
	switch name.Value {
	case PING:
		return m.Ping()
	case ECHO:
		message, ok := command.Elements[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		return m.Echo(message)
	case SET:
		key, ok := command.Elements[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		value, ok := command.Elements[2].(serializer.BulkString)
		if !ok {
			return nil
		}
		return m.Set(key, value)
	case GET:
		key, ok := command.Elements[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		return m.Get(key)
	}
	return nil
}

func (m *Manager) Ping() serializer.RESPType {
	m.logger.Info("handle command", "command", PING)
	return serializer.SimpleString("PONG")
}

func (m *Manager) Echo(message serializer.BulkString) serializer.RESPType {
	m.logger.Info("handle command", "command", ECHO)
	return message
}

func (m *Manager) Set(key, value serializer.BulkString) serializer.RESPType {
	m.logger.Info("handle command", "command", SET)
	err := m.repo.Set(key.Value, value.Value)
	if err != nil {
		return nil
	}
	return serializer.SimpleString("OK")
}

func (m *Manager) Get(key serializer.BulkString) serializer.RESPType {
	m.logger.Info("handle command", "command", GET)
	value, err := m.repo.Get(key.Value)
	if err != nil {
		m.logger.Error("failed to get value from repo", "key", key.Value)
		return serializer.BulkString{Length: -1}
	}
	return serializer.BulkString{Length: int64(len(value)), Value: value}
}
