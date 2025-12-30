package middleware

import (
	"context"
	"log"
)

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint64)
	log.Printf("%v %v", userID, ok)
	return userID, ok
}
