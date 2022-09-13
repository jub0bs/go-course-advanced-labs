package internal

import "unicode/utf8"

func IsLongEnough(username string, minLen int) bool {
	return utf8.RuneCountInString(username) >= minLen
}
