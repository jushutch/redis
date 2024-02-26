package manager

import (
	"context"

	"github.com/jushutch/redis/logging"
	"github.com/jushutch/redis/serializer"
)

func (m *Manager) handleGet(ctx context.Context, command serializer.Array) serializer.RESPType {
	key, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	value, err := m.repo.Get(ctx, key.Value)
	if err != nil {
		m.logger.With(logging.FieldsFromContext(ctx)...).Warn("failed to get value from repo", "key", key.Value, "reason", err)
		return serializer.BulkString{Length: -1, Value: ""}
	}
	return serializer.BulkString{Length: int64(len(value)), Value: value}
}
