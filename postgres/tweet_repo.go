package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	twitter "github.com/trenchesdeveloper/tweeter"
)

type TweetRepo struct {
	DB *DB
}

func NewTweetRepo(DB *DB) *TweetRepo {
	return &TweetRepo{DB: DB}
}

func (t *TweetRepo) All(ctx context.Context) ([]twitter.Tweet, error) {

	query := `SELECT * FROM tweets ORDER BY created_at DESC;`

	tweets := []twitter.Tweet{}

	if err := pgxscan.Select(ctx, t.DB.Pool, &tweets, query); err != nil {
		return nil, fmt.Errorf("error getting all tweets: %w", err)
	}

	return tweets, nil
}

func (t *TweetRepo) Create(ctx context.Context, tweet twitter.Tweet) (twitter.Tweet, error) {
	tx, err := t.DB.Pool.Begin(ctx)

	if err != nil {
		return twitter.Tweet{}, fmt.Errorf("error starting transaction: %w", err)

	}

	defer tx.Rollback(ctx)

	tweet, err = createTweet(ctx, tx, tweet)

	if err != nil {
		return twitter.Tweet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return twitter.Tweet{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return tweet, nil
}

func (t *TweetRepo) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	return getTweetByID(ctx, t.DB.Pool, id)
}

func createTweet(ctx context.Context, tx pgx.Tx, tweet twitter.Tweet) (twitter.Tweet, error) {
	query := `INSERT INTO tweets (body, user_id) VALUES ($1, $2) RETURNING *;`

	t := twitter.Tweet{}

	if err := pgxscan.Get(ctx, tx, &t, query, tweet.Body, tweet.UserID); err != nil {
		return twitter.Tweet{}, fmt.Errorf("error creating tweet: %w", err)
	}

	return t, nil
}

func getTweetByID(ctx context.Context, q pgxscan.Querier, id string) (twitter.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id = $1 LIMIT 1;`

	t := twitter.Tweet{}

	if err := pgxscan.Get(ctx, q, &t, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.Tweet{}, twitter.ErrNotFound
		}

		return twitter.Tweet{}, fmt.Errorf("error getting tweet by id: %w", err)
	}

	return t, nil
}
