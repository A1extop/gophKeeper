package errors

import "errors"

var (
	ErrUsernameAndPasswordRequired = errors.New("username and password are required")
	ErrNameLockboxRequired         = errors.New("name lockbox is required")
	ErrDataRequired                = errors.New("data is required")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrInvalidCredentials          = errors.New("invalid username or password")
	ErrUsernameTooShort            = errors.New("username must be at least 3 characters")
	ErrPasswordTooShort            = errors.New("password must be at least 6 characters")
	ErrNotFound                    = errors.New("not found")
	ErrDelete                      = errors.New("delete failed")
	ErrExists                      = errors.New("lockbox with this name already exists")
	ErrIncorrectUsername           = errors.New("incorrect login")
	ErrIncorrectPassword           = errors.New("incorrect password")
	ErrNodataToUpdate              = errors.New("no data update")
	ErrInvalidToken                = errors.New("invalid token")
	ErrInvalidTokenClaims          = errors.New("invalid token claims")
	ErrUserIdNotFoundInToken       = errors.New("user id not found in token")
	ErrLockboxNameTakenByUser      = errors.New("unique_lockbox_name_per_user")
)
