package graph

import (
	"context"
	"fmt"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/graph/models"
)

func mapUserToGQL(user twitter.User) *models.User {
	return &models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (q *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}
