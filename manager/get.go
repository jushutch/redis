package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handleGet(command serializer.Array) serializer.RESPType {
	m.logger.Info("handle command", "command", GET)
	key, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	value, err := m.repo.Get(key.Value)
	if err != nil {
		m.logger.Warn("failed to get value from repo", "key", key.Value, "reason", err)
		return serializer.BulkString{Length: -1, Value: ""}
	}
	return serializer.BulkString{Length: int64(len(value)), Value: value}
}
