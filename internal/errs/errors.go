package errs

import "errors"

var (
	ErrNilUser         = errors.New("nil user")
	ErrInvalidTGID     = errors.New("invalid tg id")
	ErrInvalidID       = errors.New("invalid user id")
	ErrInvalidLanguage = errors.New("invalid language")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPhone    = errors.New("invalid phone")
	ErrInvalidUserName = errors.New("invalid username")
)
