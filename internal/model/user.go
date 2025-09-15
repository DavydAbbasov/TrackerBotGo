package model

import "time"

type User struct {
	ID          int64
	TgUserID    int64 // 0
	UserName    *string
	PhoneNumber *string
	Email       *string
	Language    *string // nil
	TimeZone    *string
	CreatedAt   time.Time
}
