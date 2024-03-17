package jwt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/config"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var signatureType = jwa.HS256

var now = time.Now

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{
		Conf: conf,
	}
}

func (s *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (twitter.AuthToken, error) {
	token, err := jwt.ParseRequest(r, jwt.WithValidate(true), jwt.WithIssuer(s.Conf.Jwt.Issuer), jwt.WithKey(signatureType, []byte(s.Conf.Jwt.Secret)))

	if err != nil {
		return twitter.AuthToken{}, twitter.ErrInvalidAccessToken
	}

	return buildToken(token), err
}

func buildToken(token jwt.Token) twitter.AuthToken {
	return twitter.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}

}

func (s *TokenService) ParseToken(ctx context.Context, payload string) (twitter.AuthToken, error) {
	token, err := jwt.Parse([]byte(payload), jwt.WithValidate(true), jwt.WithIssuer(s.Conf.Jwt.Issuer), jwt.WithKey(signatureType, []byte(s.Conf.Jwt.Secret)))

	if err != nil {
		return twitter.AuthToken{}, twitter.ErrInvalidAccessToken
	}

	return buildToken(token), err
}

func (s *TokenService) CreateAccessToken(ctx context.Context, user twitter.User) (string, error) {
	b := setDefaultBuilder(user, twitter.AccessTokenLifetime, s.Conf)

	t, err := b.Build()

	if err != nil {
		return "", fmt.Errorf("error building token: %v", err)
	}

	token, err := jwt.Sign(t, jwt.WithKey(signatureType, []byte(s.Conf.Jwt.Secret)))

	log.Println("token", string(token))

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return string(token), nil
}

func (s *TokenService) CreateRefreshToken(ctx context.Context, user twitter.User, tokenID string) (string, error) {
	b := setDefaultBuilder(user, twitter.RefreshTokenLifetime, s.Conf)

	// add jwt id to builder
	t, err := b.JwtID(tokenID).Build()

	if err != nil {
		return "", fmt.Errorf("error building token: %v", err)
	}

	token, err := jwt.Sign(t, jwt.WithKey(signatureType, []byte(s.Conf.Jwt.Secret)))

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return string(token), nil

}

func setDefaultBuilder(user twitter.User, lifetime time.Duration, conf *config.Config) *jwt.Builder {
	builder := jwt.NewBuilder().Issuer(conf.Jwt.Issuer).Subject(user.ID).IssuedAt(now()).Expiration(now().Add(lifetime))

	return builder
}
