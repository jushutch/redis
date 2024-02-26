package manager

import (
	"context"

	"github.com/jushutch/redis/serializer"
)

func (m *Manager) handleEcho(ctx context.Context, command serializer.Array) serializer.RESPType {
	message, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	return message
}
