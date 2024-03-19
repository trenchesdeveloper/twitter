package postgres

import (
	"context"
	twitter "github.com/trenchesdeveloper/tweeter"
)

type TweetRepo struct {
	DB *DB
}

func NewTweetRepo(DB *DB) *TweetRepo {
	return &TweetRepo{DB: DB}
}

func (t TweetRepo) All(ctx context.Context) ([]twitter.Tweet, error) {
	panic("implement me")
}

func (t TweetRepo) Create(ctx context.Context, tweet twitter.Tweet) (twitter.Tweet, error) {
	panic("implement me")
}

func (t TweetRepo) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	panic("implement me")
}
