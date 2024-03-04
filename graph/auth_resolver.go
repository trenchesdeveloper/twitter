package graph

import (
	"context"
	"fmt"

	"github.com/trenchesdeveloper/tweeter/graph/models"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, registerInput *models.RegisterInput) (*models.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: Register - register"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, loginInput *models.LoginInput) (*models.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: Login - login"))
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented: Logout - logout"))
}
