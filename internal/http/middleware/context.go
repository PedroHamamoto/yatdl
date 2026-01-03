package middleware

import (
	"context"
)

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint64)
	return userID, ok
}
