package manager

import (
	"context"

	"github.com/jushutch/redis/internal/serializer"
)

func (m *Manager) handleDelete(ctx context.Context, command serializer.Array) serializer.RESPType {
	var delCount int64
	for i := 1; i < int(command.Length); i++ {
		key, ok := command.Elements[i].(serializer.BulkString)
		if !ok {
			continue
		}
		err := m.repo.Delete(ctx, key.Value)
		if err == nil {
			delCount++
		}
	}
	return serializer.Integer{Value: delCount}
}
