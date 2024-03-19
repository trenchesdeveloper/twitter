//go:build integration
// +build integration

package domain

import (
	"context"
	"github.com/stretchr/testify/require"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/faker"
	"github.com/trenchesdeveloper/tweeter/test_helper"
	"log"
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

		newTweet, err := tweetService.Create(ctx, input)

		require.NoError(t, err)
		log.Println("tweet check", newTweet.ID)
		require.NotEmpty(t, newTweet.ID, "newTweet id should be set")
		require.Equal(t, input.Body, newTweet.Body, "newTweet body should be the same")
		require.Equal(t, currentUser.ID, newTweet.UserID, "newTweet user id should be the same")
		require.NotEmpty(t, newTweet.CreatedAt, "newTweet created at should be set")
	})
}
