package xec

import (
	"encoding/json"
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

func FindStringInSlice(haystack []string, needle string) bool {
	var found bool
	for _, h := range haystack {
		if h == needle {
			found = true
		}
	}
	return found
}

func ParseConfig(C interface{}) {
	CJSON, err := json.MarshalIndent(C, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	o.Debug(fmt.Sprintf("Config in indented JSON:\n %s\n", string(CJSON)))
}
