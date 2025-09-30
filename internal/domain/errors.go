package domain

import "errors"

// (опц.) доменные ошибки
var (
	ErrOpenSessionExists = errors.New("open session already exists")
	ErrNoOpenSession     = errors.New("no open session")
	ErrInvalidTimeRange  = errors.New("end time before start")
	ErrUnknownSource     = errors.New("unknown session source")
	ErrActivityExists    = errors.New("activity already exists")
	ErrNoActivity        = errors.New("no activities")
	ErrNoExistActivity   = errors.New("	the activity does not belong to the user or does not exist")
)
