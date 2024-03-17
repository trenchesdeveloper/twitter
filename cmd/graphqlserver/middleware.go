package main

import (
	twitter "github.com/trenchesdeveloper/tweeter"
	"net/http"
)

func authMiddleware(authTokenService twitter.AuthTokenService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := authTokenService.ParseTokenFromRequest(ctx, r)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			// add user id to context
			ctx = twitter.PutUserIDIntoContext(ctx, token.Sub)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
