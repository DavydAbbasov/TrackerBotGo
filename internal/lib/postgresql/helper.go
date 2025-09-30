package helper

import (
	"strconv"
	"strings"
)

func prepareIN(s []int) string {
	strs := make([]string, len(s))
	for k, i := range s {
		strs[k] = strconv.Itoa(i)
	}
	return strings.Join(strs, ",")
}
func ToNullable(p *string) any {
	if p == nil || *p == "" {
		return nil
	}
	return *p
}
func PtrIfNotEmpty(s string) *string {
    s = strings.TrimSpace(s)
    if s == "" { return nil }
    return &s
}