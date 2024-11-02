package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotSaved       = errors.New("user not saved")
	ErrInternal           = errors.New("internal error")
)
