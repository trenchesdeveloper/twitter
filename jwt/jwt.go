package jwt

import (
	"context"
	"fmt"
	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/config"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var signatureType = jwa.RS256

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{
		Conf: conf,
	}
}

func (s *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (twitter.AuthToken, error) {
	token, err := jwt.ParseRequest(r, jwt.WithValidate(true), jwt.WithIssuer(s.Conf.Jwt.Issuer), jwt.WithKey(signatureType, s.Conf.Jwt.Secret))

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
	token, err := jwt.Parse([]byte(payload), jwt.WithValidate(true), jwt.WithIssuer(s.Conf.Jwt.Issuer), jwt.WithKey(signatureType, s.Conf.Jwt.Secret))

	if err != nil {
		return twitter.AuthToken{}, twitter.ErrInvalidAccessToken
	}

	return buildToken(token), err
}

func (s *TokenService) CreateAccessToken(ctx context.Context, user twitter.User, tokenID string) (string, error) {
	t := jwt.New()

	if err := setDefaultToken(t, user, twitter.AccessTokenLifetime, s.Conf); err != nil {
		return "", err
	}

	token, err := jwt.Sign(t, jwt.WithKey(signatureType, s.Conf.Jwt.Secret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return string(token), nil
}

func (s *TokenService) CreateRefreshToken(ctx context.Context, user twitter.User) (string, error) {
	t := jwt.New()

	if err := setDefaultToken(t, user, twitter.RefreshTokenLifetime, s.Conf); err != nil {
		return "", fmt.Errorf("error setting default token: %v", err)
	}

	if err := t.Set(jwt.JwtIDKey, user.ID); err != nil {
		return "", fmt.Errorf("error setting jwt id: %v", err)
	}

	token, err := jwt.Sign(t, jwt.WithKey(signatureType, s.Conf.Jwt.Secret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return string(token), nil

}

func setDefaultToken(t jwt.Token, user twitter.User, lifetime time.Duration, conf *config.Config) error {
	if err := t.Set(jwt.SubjectKey, user.ID); err != nil {
		return fmt.Errorf("error setting subject: %v", err)
	}

	if err := t.Set(jwt.IssuerKey, conf.Jwt.Issuer); err != nil {
		return fmt.Errorf("error setting issuer: %v", err)
	}

	if err := t.Set(jwt.IssuedAtKey, time.Now().Unix()); err != nil {
		return fmt.Errorf("error setting issued at: %v", err)
	}

	if err := t.Set(jwt.ExpirationKey, time.Now().Add(lifetime).Unix()); err != nil {
		return fmt.Errorf("error setting expiration: %v", err)
	}

	return nil
}
