//go:build integration
// +build integration

package domain

import (
	"context"
	"github.com/stretchr/testify/require"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/faker"
	"github.com/trenchesdeveloper/tweeter/test_helper"
	"testing"
)

func TestIntegrationTweetService_Create(t *testing.T) {
	t.Run("not auth user cannot create a tweet", func(t *testing.T) {
		ctx := context.Background()

		_, err := tweetService.Create(ctx, twitter.CreateTweetInput{
			Body: "a tweet",
		})

		require.ErrorIs(t, err, twitter.ErrUnAuthenicated)
	})

	t.Run("auth user can create a tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helper.TeardownDB(ctx, t, db)

		currentUser, _ := test_helper.CreateUser(ctx, t, userRepo)

		ctx = test_helper.LoginUser(ctx, t, currentUser)
		input := twitter.CreateTweetInput{
			Body: faker.RandString(100),
		}

		tweet, err := tweetService.Create(ctx, input)

		require.NoError(t, err)
		require.NotEmpty(t, tweet.ID, "tweet id should be set")
		require.Equal(t, input.Body, tweet.Body, "tweet body should be the same")
		require.Equal(t, currentUser.ID, tweet.UserID, "tweet user id should be the same")
		require.NotEmpty(t, tweet.CreatedAt, "tweet created at should be set")
	})
}
