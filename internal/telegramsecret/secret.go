// Package telegramsecret validates Telegram webhook secret tokens.
package telegramsecret

import (
	"crypto/subtle"
	stderrors "errors"
	"regexp"
)

const maxLength = 256

var validPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// Validate checks that token is empty or a valid Telegram webhook secret token.
func Validate(token string) error {
	if token == "" {
		return nil
	}
	if len(token) > maxLength {
		return stderrors.New("secret token length must be between 1 and 256")
	}
	if !validPattern.MatchString(token) {
		return stderrors.New("secret token contains invalid characters")
	}

	return nil
}

// ConstantTimeEqual compares a and b without data-dependent timing when lengths match.
func ConstantTimeEqual(a string, b string) bool {
	if len(a) != len(b) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
