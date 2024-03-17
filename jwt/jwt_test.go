package jwt

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/require"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/config"
)

var (
	conf         *config.Config
	tokenService *TokenService
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env.test")
	conf = config.New()
	tokenService = NewTokenService(conf)

	os.Exit(m.Run())
}

func TestTokenService_CreateAccessToken(t *testing.T) {
	t.Run("can create a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)

		require.NoError(t, err)

		now = func() time.Time {
			return time.Now()
		}

		require.NotEmpty(t, token)

		tok, err := jwt.Parse([]byte(token), jwt.WithKey(signatureType, []byte(conf.Jwt.Secret)), jwt.WithIssuer(conf.Jwt.Issuer))

		require.NoError(t, err)
		require.NotNil(t, tok)
		require.Equal(t, user.ID, tok.Subject())
		require.Equal(t, now().Add(twitter.AccessTokenLifetime).Unix(), tok.Expiration().Unix())
		teardownTimeNow(t)
	})

}

func TestTokenService_CreateRefreshToken(t *testing.T) {
	t.Run("can create a valid refresh token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateRefreshToken(ctx, user, "456")

		require.NoError(t, err)

		now = func() time.Time {
			return time.Now()
		}

		require.NotEmpty(t, token)

		tok, err := jwt.Parse([]byte(token), jwt.WithKey(signatureType, []byte(conf.Jwt.Secret)), jwt.WithIssuer(conf.Jwt.Issuer))

		require.NoError(t, err)
		require.NotNil(t, tok)
		require.Equal(t, user.ID, tok.Subject())
		require.Equal(t, "456", tok.JwtID())
		require.Equal(t, now().Add(twitter.RefreshTokenLifetime).Unix(), tok.Expiration().Unix())

	})

}

func TestTokenService_ParseToken(t *testing.T) {
	t.Run("can parse a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)

		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token)

		require.NoError(t, err)
		require.Equal(t, user.ID, tok.Sub)
	},
	)

	t.Run("returns error when token is invalid", func(t *testing.T) {
		ctx := context.Background()
		token, err := tokenService.ParseToken(ctx, "invalid token")

		require.Error(t, err)
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)
		require.Empty(t, token)
	})

	t.Run("returns error when token is expired", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		// make token expired
		now = func() time.Time {
			return time.Now().Add(-twitter.AccessTokenLifetime * 5)
		}

		token, err := tokenService.CreateAccessToken(ctx, user)

		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token)

		require.Error(t, err)
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)
		require.Empty(t, tok)

		teardownTimeNow(t)
	})
}

func TestTokenService_ParseTokenFromRequest(t *testing.T) {
	t.Run("can parse a valid access token from request", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		req := httptest.NewRequest("GET", "/test", nil)

		token, err := tokenService.CreateAccessToken(ctx, user)

		require.NoError(t, err)

		req.Header.Set("Authorization", token)

		tok, err := tokenService.ParseTokenFromRequest(ctx, req)

		require.NoError(t, err)
		require.Equal(t, user.ID, tok.Sub)

		// pass bearer token
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		tok, err = tokenService.ParseTokenFromRequest(ctx, req)

		require.NoError(t, err)
		require.Equal(t, user.ID, tok.Sub)


	},
	)

	t.Run("returns error when token is expired", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		req := httptest.NewRequest("GET", "/test", nil)

		// make token expired
		now = func() time.Time {
			return time.Now().Add(-twitter.AccessTokenLifetime * 5)
		}

		token, err := tokenService.CreateAccessToken(ctx, user)

		require.NoError(t, err)

		req.Header.Set("Authorization", token)

		tok, err := tokenService.ParseTokenFromRequest(ctx, req)

		require.Error(t, err)
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)
		require.Empty(t, tok)

		teardownTimeNow(t)
	})


}


func teardownTimeNow(t *testing.T) {
	t.Helper()

	now = func() time.Time {
		return time.Now()
	}
}