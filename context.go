package twitter

import (
	"context"
)

type contextKey string

var userIDKey contextKey = "userIDKey"

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID := ctx.Value(userIDKey)
	if userID == nil {
		return "", ErrNoUserIdInContext
	}

	userIDString, ok := userID.(string)
	if !ok {
		return "", ErrNoUserIdInContext
	}

	return userIDString, nil
}

func PutUserIDIntoContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
