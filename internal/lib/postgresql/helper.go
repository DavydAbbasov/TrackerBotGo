package postgresql

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
