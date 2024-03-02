package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/mocks"
)

func TestAuthService_Register(t *testing.T) {
	validInput := twitter.RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("valid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twitter.User{
			ID:        "user_id",
			Username:  validInput.Username,
			Email:     validInput.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		service := NewAuthService(userRepo)

		res, err := service.Register(ctx, validInput)

		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)

		userRepo.AssertExpectations(t)
	})
}
