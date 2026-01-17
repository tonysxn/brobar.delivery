package errors

import "errors"

var (
	UserNotFound      = errors.New("user not found")
	UserInvalidData   = errors.New("invalid user data")
	UserAlreadyExists = errors.New("user already exists")
)
