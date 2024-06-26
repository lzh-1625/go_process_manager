package utils

import (
	"regexp"
	"unicode/utf8"
)

func RemoveNotValidUtf8InString(s string) string {
	ret := s
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		ret = string(v)
	}
	return ret
}

func RemoveANSI(input string) string {
	// Define the regular expression to match ANSI escape sequences
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	// Replace all ANSI escape sequences with an empty string
	cleanedString := re.ReplaceAllString(input, "")
	return cleanedString
}
