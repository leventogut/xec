package xec

import (
	"log"
	"regexp"
)

// CheckRegex checks if the given string is a match to given pattern.
func CheckRegex(stringInput, pattern string) bool {

	m, err := regexp.MatchString(pattern, stringInput)
	if err != nil {
		log.Fatalf("Can't decode config, error: %v", err)
	}
	if m {
		return true
	} else {
		return false
	}
}
