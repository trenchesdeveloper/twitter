package twitter

import "errors"

var (
	ErrBadCredentials     = errors.New("email/password wrong combination")
	ErrValidation         = errors.New("validation error")
	ErrNotFound           = errors.New("not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrNoUserIdInContext  = errors.New("no user id in context")
	ErrGenAccessToken     = errors.New("error generating access token")
	ErrUnAuthenicated     = errors.New("unauthenticated")
)
