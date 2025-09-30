package profile

import (
	"fmt"

	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProfileView struct {
	ID, StreakDays                                    int64
	Name, Lang, TZ, Sub, PhoneNumber, Activity, Email string
}

func BuildProfileView(u *model.User, tg *tgbotapi.User) ProfileView {

	v := ProfileView{Name: "—", Lang: "—", TZ: "—", Sub: "—", PhoneNumber:"—", Activity: "—", Email: "—"}

	if tg != nil {
		v.ID = tg.ID
	}

	if u != nil {
		if u.TgUserID != 0 {
			v.ID = u.TgUserID
		}
		if u.UserName != nil && *u.UserName != "" {
			v.Name = "@" + *u.UserName
		}
		if u.Language != nil && *u.Language != "" {
			v.Lang = *u.Language
		}
		if u.TimeZone != nil && *u.TimeZone != "" {
			v.TZ = *u.TimeZone
		}
		if u.PhoneNumber != nil && *u.PhoneNumber != "" {
			v.PhoneNumber = *u.PhoneNumber
		}
		if u.Email != nil && *u.Email != "" {
			v.Email = *u.Email
		}
	}

	if v.Name == "—" && tg != nil && tg.UserName != "" {
		v.Name = "@" + tg.UserName
	}
	return v
}
func (v ProfileView) Markdown() string {
	return fmt.Sprintf(
		`
— 🛜 *%d*
— 👤 Имя: *%s*
— 🔥 Streak: *%d*
— 🌐 Язык: *%s*
— 📍 Часовой пояс: *%s*
— 🗃 Подписка: *%s*
— 📧 Контакт: %s
— 📧 Email : %s

Настроить профиль можно ниже:`,
		v.ID,
		v.Name,
		v.StreakDays,
		v.Lang,
		v.TZ,
		v.Sub,
		v.PhoneNumber,
		v.Email,
	)
}
