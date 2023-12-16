package manager

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jushutch/redis/serializer"
)

type ExpirationOpt string

const (
	EX   ExpirationOpt = "EX"
	PX   ExpirationOpt = "PX"
	EXAT ExpirationOpt = "EXAT"
	PXAT ExpirationOpt = "PXAT"
)

func IsExpirationOpt(opt string) bool {
	o := ExpirationOpt(opt)
	return o == EX || o == PX || o == EXAT || o == PXAT
}

func (o ExpirationOpt) GetExpiration(raw string) (int64, error) {
	switch o {
	case EX:
		expiration, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return 0, nil
		}
		return time.Now().Add(time.Duration(expiration) * time.Second).UnixMilli(), nil
	case PX:
		expiration, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return 0, nil
		}
		return time.Now().Add(time.Duration(expiration) * time.Millisecond).UnixMilli(), nil
	case EXAT:
		expiration, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return 0, nil
		}
		return time.Unix(expiration, 0).UnixMilli(), nil
	case PXAT:
		return strconv.ParseInt(raw, 10, 0)
	default:
		return 0, fmt.Errorf("%q is not an expiration option", o)
	}
}

func (m *Manager) handleSet(command serializer.Array) serializer.RESPType {
	m.logger.Info("handle command", "command", SET)

	key, ok := command.Elements[1].(serializer.BulkString)
	if !ok {
		return nil
	}
	value, ok := command.Elements[2].(serializer.BulkString)
	if !ok {
		return nil
	}
	var expirationUnix int64 = 0
	var err error
	if command.Length == 5 {
		option, ok := command.Elements[3].(serializer.BulkString)
		if !ok {
			return nil
		}
		expirationStr, ok := command.Elements[4].(serializer.BulkString)
		if !ok {
			return nil
		}
		expirationUnix, err = ExpirationOpt(option.Value).GetExpiration(expirationStr.Value)
		if err != nil {
			return nil
		}
	}
	err = m.repo.Set(key.Value, value.Value, expirationUnix)
	if err != nil {
		return nil
	}
	return serializer.SimpleString("OK")
}
