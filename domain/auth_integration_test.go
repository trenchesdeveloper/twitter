//go:build integration
// +build integration

package domain

import (
	"context"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/faker"
	"github.com/trenchesdeveloper/tweeter/test_helper"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := twitter.RegisterInput{

		Username:        faker.Username(),
		Email:           faker.Email(),
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("can register user", func(t *testing.T) {
		ctx := context.Background()

		defer test_helper.TeardownDB(ctx, t, db)

		res, err := authService.Register(ctx, validInput)

		require.NoError(t, err)

		require.NotEmpty(t, res.User.ID)

		require.Equal(t, validInput.Username, res.User.Username)

		require.Equal(t, validInput.Email, res.User.Email)

		require.NotEqual(t, validInput.Password, res.User.Password)

	})

	t.Run("Existing username", func(t *testing.T) {
		ctx := context.Background()

		defer test_helper.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)

		require.NoError(t, err)

		_, err = authService.Register(ctx, twitter.RegisterInput{
			Username:        validInput.Username,
			Email:           "new@gmail.com",
			Password:        "password",
			ConfirmPassword: "password",
		})

		require.Error(t, err)

		require.ErrorIs(t, err, twitter.ErrUsernameTaken)

	})

	t.Run("Existing email", func(t *testing.T) {
		ctx := context.Background()

		defer test_helper.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)

		require.NoError(t, err)

		_, err = authService.Register(ctx, twitter.RegisterInput{
			Username:        "new",
			Email:           validInput.Email,
			Password:        "password",
			ConfirmPassword: "password",
		})

		require.Error(t, err)

		require.ErrorIs(t, err, twitter.ErrEmailTaken)

	})

}
