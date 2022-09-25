package chat

import (
	"unicode"
)

type user struct {
	Name string
}

func (u user) String() string {
	return u.Name
}

func (u user) IsValid() bool {
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
