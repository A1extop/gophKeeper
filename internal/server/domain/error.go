package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUserType    = errors.New("invalid user type")
	ErrNameEmpty          = errors.New("name is empty")
	ErrNoDataToCreate     = errors.New("no data to create")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrTokenCreation = errors.New("could not create token")
)

var (
	ErrLockBoxNotFound = errors.New("lockbox not found")
)
var (
	ErrInvalidUserID = errors.New("invalid user ID")
)
