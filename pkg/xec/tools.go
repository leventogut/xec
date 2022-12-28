package xec

import (
	"fmt"
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

// ConvertEnvMapToEnvString converts a map[string]string into key=value format.
func ConvertEnvMapToEnvString(kvm map[string]string) string {
	for k, v := range kvm {
		return fmt.Sprintf("%s=%s", k, v)
	}
	return ""
}
