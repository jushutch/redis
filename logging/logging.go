package logging

import (
	"context"

	"github.com/jushutch/redis/tracing"
)

func FieldsFromContext(ctx context.Context) []any {
	if traceID := tracing.GetTraceID(ctx); traceID != "" {
		return []any{"trace-id", traceID}
	}

	return []any{}
}
