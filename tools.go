package main

import (
	"strings"
	"unicode"
)

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func checkIfStringExists(strs, str string) bool {
	elements := strings.Split(strs, ",")
	for _, element := range elements {
		if element == str {
			return true
		}
	}
	return false
}

func trimExtraWhitespace(str string) string {
	var builder strings.Builder
	inWhitespace := false

	for _, runeValue := range str {
		if unicode.IsSpace(runeValue) {
			if !inWhitespace && builder.Len() > 0 {
				builder.WriteRune(' ')
			}
			inWhitespace = true
		} else {
			builder.WriteRune(runeValue)
			inWhitespace = false
		}
	}
	result := builder.String()
	if inWhitespace && len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}
