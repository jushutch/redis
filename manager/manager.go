package manager

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jushutch/redis/serializer"
)

type Command string

// Define supported commands
const (
	PING      Command = "PING"
	ECHO      Command = "ECHO"
	SET       Command = "SET"
	GET       Command = "GET"
	EXISTS    Command = "EXISTS"
	DELETE    Command = "DEL"
	INCREMENT Command = "INCR"
)

// Repo manages the storing and retreiving of data
type Repo interface {
	Set(key string, value string, expiration int64) error
	Get(key string) (string, error)
	Delete(key string) error
	Increment(key string) (int64, error)
}

// Manager handles Redis commands
type Manager struct {
	repo   Repo
	logger *slog.Logger
}

// New creates a new Redis manager
func New(repo Repo, logger *slog.Logger) *Manager {
	return &Manager{
		logger: logger.With("name", "redis.manager"),
		repo:   repo,
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
	switch Command(strings.ToUpper(name.Value)) {
	case PING:
		return m.handlePing(command)
	case ECHO:
		return m.handleEcho(command)
	case SET:
		return m.handleSet(command)
	case GET:
		return m.handleGet(command)
	case EXISTS:
		return m.handleExists(command)
	case DELETE:
		return m.handleDelete(command)
	case INCREMENT:
		return m.handleIncrement(command)
	default:
		m.logger.Warn("unknown command", "command", name.Value)
		return serializer.SimpleError(fmt.Sprintf("ERR unknown command '%s'", name.Value))
	}
}
