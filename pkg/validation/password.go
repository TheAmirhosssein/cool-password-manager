package validation

import "regexp"

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#\$%\^&\*]`)

	return upper.MatchString(password) &&
		lower.MatchString(password) &&
		number.MatchString(password) &&
		special.MatchString(password)
}
