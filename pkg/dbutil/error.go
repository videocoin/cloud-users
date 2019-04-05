package dbutil

import (
	"strconv"
	"unicode"
)

var (
	ErrDuplicateEntry = 1062
)

func ErrorCode(e error) int {
	err := e.Error()

	if len(err) < 6 {
		return 0
	}
	i := 6

	for ; len(err) > i && unicode.IsDigit(rune(err[i])); i++ {
	}

	n, e := strconv.Atoi(string(err[6:i]))
	if e != nil {
		return 0
	}

	return n
}
