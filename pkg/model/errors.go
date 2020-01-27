package model

import "errors"

var (
	// ErrUserNotProvided occurs when user not found in context
	ErrUserNotProvided = errors.New("user not provided")
)
