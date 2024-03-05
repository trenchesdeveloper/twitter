package graph

import (
	"context"
	"errors"
	"fmt"
	twitter "github.com/trenchesdeveloper/tweeter"

	"github.com/trenchesdeveloper/tweeter/graph/models"
)

func mapAuthResponseToGQL(res twitter.AuthResponse) *models.AuthResponse {
	return &models.AuthResponse{
		User:        mapUserToGQL(res.User),
		AccessToken: res.AccessToken,
	}
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, registerInput *models.RegisterInput) (*models.AuthResponse, error) {
	res, err := r.AuthService.Register(ctx, twitter.RegisterInput{
		Username:        registerInput.Username,
		Email:           registerInput.Email,
		Password:        registerInput.Password,
		ConfirmPassword: registerInput.ConfirmPassword,
	})

	if err != nil {
		switch {
		case errors.Is(err, twitter.ErrValidation) ||
			errors.Is(err, twitter.ErrUsernameTaken) ||
			errors.Is(err, twitter.ErrEmailTaken):
			return nil, buildBadRequestError(ctx, err)

		default:
			return nil, err
		}

	}

	return mapAuthResponseToGQL(res), nil

}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, loginInput *models.LoginInput) (*models.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: Login - login"))
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented: Logout - logout"))
}
