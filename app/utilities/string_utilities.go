package utilities

import (
	"strings"
	"unicode"
	"regexp"
)

var removeAlphNumRegex = regexp.MustCompile(`[^\p{L}\p{N} ]+`)

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

func ConvertToFilename(s string) string {
	nonSpecialStr := removeAlphNumRegex.ReplaceAllString(s, " ")
	lowerStr := strings.ToLower(nonSpecialStr)
	return strings.Join(strings.Fields(lowerStr), "_")
}
