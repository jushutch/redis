package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handlePing(_ serializer.Array) serializer.RESPType {
	return serializer.SimpleString("PONG")
}
