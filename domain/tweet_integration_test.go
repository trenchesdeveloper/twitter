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

func TestIntegrationTweetService_All(t *testing.T) {
	t.Run("can get all tweets", func(t *testing.T) {
		ctx := context.Background()
		defer test_helper.TeardownDB(ctx, t, db)

		user, _ := test_helper.CreateUser(ctx, t, userRepo)

		test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweets, err := tweetService.All(ctx)

		require.NoError(t, err)
		require.NotEmpty(t, tweets)
		require.Len(t, tweets, 3)

	})
}

func TestIntegrationTweetService_GetByID(t *testing.T) {
	t.Run("can get tweet by id", func(t *testing.T) {
		ctx := context.Background()
		defer test_helper.TeardownDB(ctx, t, db)

		user, _ := test_helper.CreateUser(ctx, t, userRepo)

		existingTweet := test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweet, err := tweetService.GetByID(ctx, existingTweet.ID)

		require.NoError(t, err)
		require.Equal(t, existingTweet.ID, tweet.ID)
		require.Equal(t, existingTweet.Body, tweet.Body)
		require.Equal(t, existingTweet.UserID, tweet.UserID)
		require.Equal(t, existingTweet.CreatedAt, tweet.CreatedAt)
	})

	t.Run("cannot get tweet by id if tweet does not exist", func(t *testing.T) {
		ctx := context.Background()

		_, err := tweetService.GetByID(ctx, faker.UUID())

		require.ErrorIs(t, err, twitter.ErrNotFound)
	})

	t.Run("cannot get tweet by id if id is invalid", func(t *testing.T) {
		ctx := context.Background()

		_, err := tweetService.GetByID(ctx, "invalid-uuid")

		require.ErrorIs(t, err, twitter.ErrInvalidUUID)
	})
}
