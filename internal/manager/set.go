package manager

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jushutch/redis/internal/logging"
	"github.com/jushutch/redis/internal/serializer"
)

type ExpirationOpt string

const (
	EX   ExpirationOpt = "EX"
	PX   ExpirationOpt = "PX"
	EXAT ExpirationOpt = "EXAT"
	PXAT ExpirationOpt = "PXAT"
)

func (o ExpirationOpt) GetExpiration(raw string) (int64, error) {
	switch o {
	case EX:
		expiration, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return 0, err
		}
		return time.Now().Add(time.Duration(expiration) * time.Second).UnixMilli(), nil
	case PX:
		expiration, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return 0, err
		}
		return time.Now().Add(time.Duration(expiration) * time.Millisecond).UnixMilli(), nil
	case EXAT:
		expiration, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return 0, err
		}
		return time.Unix(expiration, 0).UnixMilli(), nil
	case PXAT:
		return strconv.ParseInt(raw, 10, 0)
	default:
		return 0, fmt.Errorf("%q is not an expiration option", o)
	}
}

func (m *Manager) handleSet(ctx context.Context, command serializer.Array) serializer.RESPType {
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
		expirationUnix, err = ExpirationOpt(strings.ToUpper(option.Value)).GetExpiration(expirationStr.Value)
		if err != nil {
			m.logger.With(logging.FieldsFromContext(ctx)...).Error("failed to parse expiration argument", "error", err)
			return nil
		}
	}
	err = m.repo.Set(ctx, key.Value, value.Value, expirationUnix)
	if err != nil {
		return serializer.SimpleError(fmt.Sprintf("ERR %s", err.Error()))
	}
	return serializer.SimpleString("OK")
}
