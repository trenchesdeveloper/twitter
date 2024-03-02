package domain

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	twitter "github.com/trenchesdeveloper/tweeter"
)

type AuthService struct {
	UserRepo twitter.UserRepo
}

func NewAuthService(userRepo twitter.UserRepo) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, input twitter.RegisterInput) (twitter.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.AuthResponse{}, err
	}

	// check if user exists
	if _, err := s.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrUsernameTaken
	}

	// check if email exists
	if _, err := s.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrEmailTaken
	}

	user := twitter.User{
		Username: input.Username,
		Email:    input.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("error hashing password %v", err)
	}

	user.Password = string(hashedPassword)

	user, err = s.UserRepo.Create(ctx, user);

	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("error creating user %v", err)
	}

	// token, err := generateToken(user.ID)
	// if err != nil {
	// 	return twitter.AuthResponse{}, fmt.Errorf("%w: error generating token", twitter.ErrInternal)
	// }

	return twitter.AuthResponse{
		AccessToken: "token",
		User:        user,
	}, nil
}


func (s *AuthService) Login(ctx context.Context, input twitter.LoginInput) (twitter.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.AuthResponse{}, err
	}

	user, err := s.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, twitter.ErrNotFound):
			return twitter.AuthResponse{}, twitter.ErrBadCredentials
		default:
			return twitter.AuthResponse{}, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return twitter.AuthResponse{}, twitter.ErrBadCredentials
	}

	// token, err := generateToken(user.ID)
	// if err != nil {
	// 	return twitter.AuthResponse{}, fmt.Errorf("%w: error generating token", twitter.ErrInternal)
	// }

	return twitter.AuthResponse{
		AccessToken: "token",
		User:        user,
	}, nil
}