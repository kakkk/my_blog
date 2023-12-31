package utils

import (
	"net/url"
	"regexp"
)

func CheckIsEmail(str string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return reg.MatchString(str)
}

func CheckIsURL(str string) bool {
	if str == "" {
		return true
	}
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
