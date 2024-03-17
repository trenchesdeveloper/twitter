package jwt

import (
	"context"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/require"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/config"
	"os"
	"testing"
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
		require.NotEmpty(t, token)

		tok, err := jwt.Parse([]byte(token), jwt.WithKey(signatureType, []byte(conf.Jwt.Secret)), jwt.WithIssuer(conf.Jwt.Issuer))

		require.NoError(t, err)
		require.NotNil(t, tok)
		require.Equal(t, user.ID, tok.Subject())

	})

}
