package manager

import (
	"context"

	"github.com/jushutch/redis/internal/logging"
	"github.com/jushutch/redis/internal/serializer"
)

func (m *Manager) handleIncrement(ctx context.Context, command serializer.Array) serializer.RESPType {
	key, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	newValue, err := m.repo.Add(ctx, key.Value, 1)
	if err != nil {
		m.logger.With(logging.FieldsFromContext(ctx)...).Warn("failed to increment value", "key", key, "error", err)
		return serializer.SimpleError("ERR value is not an integer or out of range")
	}
	return serializer.Integer{Value: newValue}
}
