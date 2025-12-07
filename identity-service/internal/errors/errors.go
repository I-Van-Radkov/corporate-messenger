package errors

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")

	ErrTokenExpired = errors.New("token is expired")
	ErrInvalidToken = errors.New("invalid token")
)
