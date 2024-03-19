package twitter

import (
	"context"
	"strings"
	"time"
)

var (
	TweetMaxLength = 140
	TweetMinLength = 1
)

type Tweet struct {
	ID        string
	Body      string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTweetInput struct {
	Body string `json:"body"`
}

func (in *CreateTweetInput) Sanitize() {
	in.Body = strings.TrimSpace(in.Body)
}

func (in CreateTweetInput) Validate() error {
	if len(in.Body) < TweetMinLength {
		return ErrValidation
	}

	if len(in.Body) > TweetMaxLength {
		return ErrValidation
	}
	return nil
}

type TweetService interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, in CreateTweetInput) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
}

type TweetRepo interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, tweet Tweet) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
}
