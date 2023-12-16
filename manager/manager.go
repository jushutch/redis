package manager

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jushutch/redis/repo"
	"github.com/jushutch/redis/serializer"
)

// Define supported commands
const (
	PING = "PING"
	ECHO = "ECHO"
	SET  = "SET"
	GET  = "GET"
)

// Repo manages the storing and retreiving of data
type Repo interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

// Manager handles Redis commands
type Manager struct {
	repo   Repo
	logger *slog.Logger
}

// New creates a new Redis manager
func New(logger *slog.Logger) *Manager {
	return &Manager{
		logger: logger.With("name", "redis.manager"),
		repo:   repo.New(logger),
	}
}

// HandleCommand executes the given command and returns a RESP response
func (m *Manager) HandleCommand(command serializer.Array) serializer.RESPType {
	if command.Length <= 0 {
		return nil
	}

	name, ok := command.Elements[0].(serializer.BulkString)
	if !ok {
		return nil
	}
	switch strings.ToUpper(name.Value) {
	case PING:
		return m.ping()
	case ECHO:
		message, ok := command.Elements[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		return m.echo(message)
	case SET:
		key, ok := command.Elements[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		value, ok := command.Elements[2].(serializer.BulkString)
		if !ok {
			return nil
		}
		return m.set(key, value)
	case GET:
		key, ok := command.Elements[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		return m.get(key)
	default:
		m.logger.Warn("unknown command", "command", name.Value)
		return serializer.SimpleError(fmt.Sprintf("ERR unknown command '%s'", name.Value))
	}
}

func (m *Manager) ping() serializer.RESPType {
	m.logger.Info("handle command", "command", PING)
	return serializer.SimpleString("PONG")
}

func (m *Manager) echo(message serializer.BulkString) serializer.RESPType {
	m.logger.Info("handle command", "command", ECHO)
	return message
}

func (m *Manager) set(key, value serializer.BulkString) serializer.RESPType {
	m.logger.Info("handle command", "command", SET)
	err := m.repo.Set(key.Value, value.Value)
	if err != nil {
		return nil
	}
	return serializer.SimpleString("OK")
}

func (m *Manager) get(key serializer.BulkString) serializer.RESPType {
	m.logger.Info("handle command", "command", GET)
	value, err := m.repo.Get(key.Value)
	if err != nil {
		m.logger.Error("failed to get value from repo", "key", key.Value, "error", err)
		return serializer.BulkString{Length: -1, Value: ""}
	}
	return serializer.BulkString{Length: int64(len(value)), Value: value}
}
