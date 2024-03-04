package graph

import (
	"context"
	"fmt"
	"github.com/trenchesdeveloper/tweeter/graph/models"
)

func (q *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}
