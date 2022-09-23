package chat

import (
	"unicode"
)

type User struct {
	Name string
}

func (u User) IsValid() bool {
	if len(u.Name) == 0 || !isAlphaNumeric(u.Name) {
		return false
	}
	return true
}

func isAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
