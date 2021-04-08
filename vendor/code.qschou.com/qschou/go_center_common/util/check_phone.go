package util

import (
	"regexp"
	"strings"
)

func CheckPhone(phone string, countryCode string) bool {
	if v := FetchCountryByCode(countryCode); v == nil {
		return false
	}
	pattern := "^[0-9]{5,}$"
	if strings.ToUpper(countryCode) == "CN" {
		pattern = "^[0-9]{11}$"
	}
	if b, err := regexp.Match(pattern, []byte(phone)); err != nil {
		return false
	} else if !b {
		return false
	}
	return true
}
