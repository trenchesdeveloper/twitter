package test_helper

import (
	"context"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/faker"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trenchesdeveloper/tweeter/postgres"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)

	require.NoError(t, err)
}

func CreateUser(ctx context.Context, t *testing.T, userRepo twitter.UserRepo) (twitter.User, error) {
	t.Helper()

	user, err := userRepo.Create(ctx, twitter.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password,
	})

	require.NoError(t, err)

	return user, nil
}

func LoginUser(ctx context.Context, t *testing.T, user twitter.User) context.Context {
	t.Helper()

	return twitter.PutUserIDIntoContext(ctx, user.ID)

}
