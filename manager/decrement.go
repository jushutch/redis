package manager

import (
	"github.com/jushutch/redis/serializer"
)

func (m *Manager) handleDecrement(command serializer.Array) serializer.RESPType {
	key, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	newValue, err := m.repo.Add(key.Value, -1)
	if err != nil {
		m.logger.Warn("failed to decrement value", "key", key, "error", err)
		return serializer.SimpleError("ERR value is not an integer or out of range")
	}
	return serializer.Integer{Value: newValue}
}
