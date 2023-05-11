package proxy

import (
	"strings"
	"unicode"
)

const TonysAddress = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

/*
A substring is considered to be a Boguscoin address if it satisfies all of:

    it starts with a "7"
    it consists of at least 26, and at most 35, alphanumeric characters
    it starts at the start of a chat message, or is preceded by a space
    it ends at the end of a chat message, or is followed by a space
*/

func ReplaceBogusCoins(s string) string {
	var coins []string
	for _, f := range strings.Fields(s) {
		if IsBogusCoinAddress(f) {
			coins = append(coins, f, TonysAddress)
		}
	}
	if len(coins) > 0 {
		r := strings.NewReplacer(coins...)
		return r.Replace(s)
	}
	return s
}

func IsBogusCoinAddress(s string) bool {
	// length constaint
	if len(s) < 26 || len(s) > 35 {
		return false
	}

	// begins with 7
	if s[0] != '7' {
		return false
	}

	// is alphanumeric
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}

	return true
}
