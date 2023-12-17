package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handleEcho(command serializer.Array) serializer.RESPType {
	message, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	return message
}
