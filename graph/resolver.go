package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
)

//go:generate go run ../scripts/gqlgen.go

type Resolver struct {
	AuthService twitter.AuthService
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}
