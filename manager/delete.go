package manager

import "github.com/jushutch/redis/serializer"

func (m *Manager) handleDelete(command serializer.Array) serializer.RESPType {
	var delCount int64
	for i := 1; i < int(command.Length); i++ {
		key, ok := command.Elements[i].(serializer.BulkString)
		if !ok {
			continue
		}
		err := m.repo.Delete(key.Value)
		if err == nil {
			delCount++
		}
	}
	return serializer.Integer{Value: delCount}
}
