package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handlePing(_ serializer.Array) serializer.RESPType {
	m.logger.Info("handle command", "command", PING)
	return serializer.SimpleString("PONG")
}
