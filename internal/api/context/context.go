package context

import "context"

type contextKey int

const userIDKey contextKey = iota

func WithUserID(ctx context.Context, userID uint64) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) (uint64, bool) {
	if ctx == nil {
		return 0, false
	}

	eID, ok := ctx.Value(userIDKey).(uint64)
	if eID == 0 {
		return 0, false
	}

	return eID, ok
}
