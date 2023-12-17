package manager

import (
	"github.com/jushutch/redis/serializer"
)

func (m *Manager) handleIncrement(command serializer.Array) serializer.RESPType {
	m.logger.Info("handle command", "command", INCREMENT)
	key, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	newValue, err := m.repo.Increment(key.Value)
	if err != nil {
		return serializer.SimpleError("ERR value is not an integer or out of range")
	}
	return serializer.Integer{Value: newValue}
}
