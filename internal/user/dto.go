package user
//для входящих данных от Telegram
type UserInput struct {
	TgUserID    int64
	UserName    *string
	PhoneNumber *string
	Email       *string
	Language    *string
	TimeZone    *string
}

func S(s string) *string { return &s }
