package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors   map[string]string
	Messages []string
}

func (v *Validator) Valid() bool {
	return len(v.Errors)+len(v.Messages) == 0
}

func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	// don't override existing error
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) AddMessage(message string) {
	v.Messages = append(v.Messages, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MaxChars(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

func MinChars(s string, n int) bool {
	return utf8.RuneCountInString(s) >= n
}

func Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}

func Equals(s1 string, s2 string) bool {
	return s1 == s2
}

func AllowedValue[T comparable](value T, allowed ...T) bool {
	return slices.Contains(allowed, value)
}
