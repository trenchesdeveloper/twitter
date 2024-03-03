package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	twitter "github.com/trenchesdeveloper/tweeter"
)

type UserRepo struct {
	DB *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (u *UserRepo) Create(ctx context.Context, user twitter.User) (twitter.User, error) {
	tx, err := u.DB.Pool.Begin(ctx)
	if err != nil {
		return twitter.User{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	user, err = createUser(ctx, tx, user)

	if err != nil {
		return twitter.User{}, fmt.Errorf("error creating user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return twitter.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return user, nil

}

func createUser(ctx context.Context, tx pgx.Tx, user twitter.User) (twitter.User, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *;`

	if err := pgxscan.Get(ctx, tx, &user, query, user.Username, user.Email, user.Password); err != nil {
		return twitter.User{}, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (twitter.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`

	user := twitter.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &user, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	return user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

	user := twitter.User{}

	if err := pgxscan.Get(ctx, u.DB.Pool, &user, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}
