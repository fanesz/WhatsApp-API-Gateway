package utils

import (
	"regexp"
	"strings"
)

// Result: parsing number 0 -> 62, and remove all non-numeric characters
// because the phone number format in Indonesia is 62xxxxxxxxxx
func ParsePhoneNumber(number string) string {
	re := regexp.MustCompile("[^0-9]")
	cleaned := re.ReplaceAllString(number, "")
	if strings.HasPrefix(cleaned, "0") {
		cleaned = "62" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "+") {
		return cleaned[1:]
	}

	return cleaned
}
