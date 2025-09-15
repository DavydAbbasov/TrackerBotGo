package user

import (
	"regexp"
	"strings"

	"github.com/DavydAbbasov/trecker_bot/internal/errs"
)

type UserValidator struct {
	AllowedLang map[string]struct{}
	EmailRegexp *regexp.Regexp
	// UserNameRules UserNameRules
}

func NewUserValidator() UserValidator {
	return UserValidator{
		AllowedLang: map[string]struct{}{
			"ru": {}, "en": {}, "de": {}, "uk": {}, "tur": {}, "arab": {},
		},
	}
}
func (v UserValidator) ValidateLanguage(input string) (string, error) {
	lang := strings.ToLower(strings.TrimSpace(input))
	if _, ok := v.AllowedLang[lang]; !ok {
		return "", errs.ErrInvalidLanguage
	}
	return lang, nil
}
