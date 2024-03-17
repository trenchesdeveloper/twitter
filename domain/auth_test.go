package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/faker"
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

		authTokenService := &mocks.AuthTokenService{}

		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("a token", nil)

		service := NewAuthService(userRepo, authTokenService)

		res, err := service.Register(ctx, validInput)

		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("can't generate access token", func(t *testing.T) {
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

		authTokenService := &mocks.AuthTokenService{}

		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("", errors.New("something went wrong"))

		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Register(ctx, validInput)

		require.ErrorIs(t, err, twitter.ErrGenAccessToken)



		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)

	})

}

func TestAuthService_Login(t *testing.T) {
	validInput := twitter.LoginInput{
		Email:    "bob@gmail.com",
		Password: "password",
	}

	hashedPassword := faker.Password
	t.Run("can login", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{
			Email:    validInput.Email,
			Password: hashedPassword,
		}, nil)

		authTokenService := &mocks.AuthTokenService{}

		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("a token", nil)

		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Login(ctx, validInput)

		require.NoError(t, err)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{
			Email:    validInput.Email,
			Password: hashedPassword,
		}, nil)

		authTokenService := &mocks.AuthTokenService{}

		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Login(ctx, twitter.LoginInput{
			Email:    validInput.Email,
			Password: "wrong_password",
		})

		require.ErrorIs(t, err, twitter.ErrBadCredentials)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Login(ctx, validInput)

		require.ErrorIs(t, err, twitter.ErrBadCredentials)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("get user by email error", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{},
			errors.New("something went wrong"))

		authTokenService := &mocks.AuthTokenService{}

		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Login(ctx, validInput)

		require.Error(t, err)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		authTokenService := &mocks.AuthTokenService{}

		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Login(ctx, twitter.LoginInput{
			Email:    "bob.com",
			Password: "",
		})

		require.ErrorIs(t, err, twitter.ErrValidation)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("can't generate access token", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{
			Email:    validInput.Email,
			Password: hashedPassword,
		}, nil)

		authTokenService := &mocks.AuthTokenService{}

		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("", errors.New("something went wrong"))

		service := NewAuthService(userRepo, authTokenService)

		_, err := service.Login(ctx, validInput)

		require.ErrorIs(t, err, twitter.ErrGenAccessToken)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})
}
