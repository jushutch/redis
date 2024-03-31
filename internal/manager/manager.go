package manager

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jushutch/redis/internal/logging"
	"github.com/jushutch/redis/internal/serializer"
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
	DECREMENT Command = "DECR"
)

// Repo manages the storing and retreiving of data
type Repo interface {
	Set(ctx context.Context, key string, value string, expiration int64) error
	Get(ctx context.Context, key string) (string, error)
	Add(ctx context.Context, key string, toAdd int64) (int64, error)
	Delete(ctx context.Context, key string) error
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
func (m *Manager) HandleCommand(ctx context.Context, command serializer.Array) serializer.RESPType {
	if command.Length <= 0 {
		return nil
	}

	name, ok := command.Elements[0].(serializer.BulkString)
	if !ok {
		return nil
	}
	commandName := Command(strings.ToUpper(name.Value))
	m.logger.With(logging.FieldsFromContext(ctx)...).Info("handle command", "command", commandName)
	switch commandName {
	case PING:
		return m.handlePing(ctx, command)
	case ECHO:
		return m.handleEcho(ctx, command)
	case SET:
		return m.handleSet(ctx, command)
	case GET:
		return m.handleGet(ctx, command)
	case EXISTS:
		return m.handleExists(ctx, command)
	case DELETE:
		return m.handleDelete(ctx, command)
	case INCREMENT:
		return m.handleIncrement(ctx, command)
	case DECREMENT:
		return m.handleDecrement(ctx, command)
	default:
		return serializer.SimpleError(fmt.Sprintf("ERR unknown command '%s'", name.Value))
	}
}
