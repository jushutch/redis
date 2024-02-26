package manager

import (
	"context"

	"github.com/jushutch/redis/serializer"
)

func (m *Manager) handlePing(ctx context.Context, _ serializer.Array) serializer.RESPType {
	return serializer.SimpleString("PONG")
}
