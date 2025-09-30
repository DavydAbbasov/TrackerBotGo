package helper

import (
	"strings"

	"github.com/DavydAbbasov/trecker_bot/internal/errs"
)

type Validator struct {
	AllowedLangs map[string]struct{}
	// EmailRegexp *regexp.Regexp
	// UserNameRules UserNameRules
}

func NewUserValidator() *Validator {
	return &Validator{
		AllowedLangs: map[string]struct{}{
			"ru": {}, "en": {}, "de": {}, "uk": {}, "ar": {},
		},
	}
}
func NormalizeLang(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "_", "-")
	if i := strings.IndexByte(s, '-'); i > 0 { // en-US -> en
		s = s[:i]
	}
	if s == "ua" {
		s = "uk"
	} // синоним
	return s
}
func (v *Validator) ValidateLanguage(raw string) (string, error) {
	code := NormalizeLang(raw)
	if _, ok := v.AllowedLangs[code]; !ok {
		return "", errs.ErrInvalidLanguage
	}
	return code, nil
}
