package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handleEcho(command serializer.Array) serializer.RESPType {
	m.logger.Info("handle command", "command", ECHO)
	message, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	return message
}
