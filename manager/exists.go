package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handleExists(command serializer.Array) serializer.RESPType {
	var existCount int64
	for i := 1; i < int(command.Length); i++ {
		key, ok := command.Elements[i].(serializer.BulkString)
		if !ok {
			continue
		}
		_, err := m.repo.Get(key.Value)
		if err == nil {
			existCount++
		}
	}
	return serializer.Integer{Value: existCount}
}
