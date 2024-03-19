package domain

import (
	"context"
	twitter "github.com/trenchesdeveloper/tweeter"
)

type TweetService struct {
	TweetRepo twitter.TweetRepo
}

func NewTweetService(tr twitter.TweetRepo) *TweetService {
	return &TweetService{TweetRepo: tr}
}

func (t TweetService) All(ctx context.Context) ([]twitter.Tweet, error) {
	// check if user is authenticated
	userID, err := twitter.GetUserIDFromContext(ctx)

	if err != nil {
		return nil, err
	}

	return []twitter.Tweet{
		{
			ID: userID,
		},
	}, nil
}

func (t TweetService) Create(ctx context.Context, in twitter.CreateTweetInput) (twitter.Tweet, error) {
	// check if user is authenticated
	_, err := twitter.GetUserIDFromContext(ctx)

	if err != nil {
		return twitter.Tweet{}, twitter.ErrUnAuthenicated

	}

	return twitter.Tweet{}, nil
}

func (t TweetService) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	//TODO implement me
	panic("implement me")
}