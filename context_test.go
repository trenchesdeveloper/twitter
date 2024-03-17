package twitter

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("get user id from context", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, userIDKey, "123")

		userID, err := GetUserIDFromContext(ctx)

		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})

	t.Run("return error when user id is not in context", func(t *testing.T) {
		ctx := context.Background()

		userID, err := GetUserIDFromContext(ctx)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrNoUserIdInContext)
		require.Empty(t, userID)
	})

	t.Run("return error when user id is not a string", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, userIDKey, 123)

		userID, err := GetUserIDFromContext(ctx)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrNoUserIdInContext)
		require.Empty(t, userID)
	})

}

func TestPutUserIDIntoContext(t *testing.T) {
	t.Run("put user id into context", func(t *testing.T) {
		ctx := context.Background()

		ctx = PutUserIDIntoContext(ctx, "123")

		userID, err := GetUserIDFromContext(ctx)

		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})
}
