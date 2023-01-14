package utilities

import (
	"strings"
	"unicode"
)

func ContainsWord(str, word string) bool {
	strArr := strings.FieldsFunc(str, delims)
	for _, s := range strArr {
		if s == word {
			return true
		}
	}

	return false
}

func delims(r rune) bool {
	return (unicode.IsPunct(r) || unicode.IsSpace(r)) && r != rune('_')
}
