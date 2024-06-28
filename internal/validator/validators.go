package Validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

//Regex to check if the email is a valid one

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	FieldErrors map[string]string
  NonFieldErrors [] string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) ==0 && len(v.NonFieldErrors) ==0
}

func (v *Validator) AddField(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exits := v.FieldErrors[key]; !exits {
		v.FieldErrors[key] = message
	}
}
func (v *Validator) AddNonFieldError(message string) {
  v.NonFieldErrors = append(v.NonFieldErrors,message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddField(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, minChr int) bool {
  return utf8.RuneCountInString(value) >=minChr
}
func Matches(value string,rx *regexp.Regexp) bool {
  return rx.MatchString(value)
}

func MaxChars(value string, maxchr int) bool {
	return utf8.RuneCountInString(value) <= maxchr
}
func PermittedValue [T comparable](value T, permitttedValues ...T) bool {
	for i := range permitttedValues {
		if value == permitttedValues[i] {
			return true
		}
	}
	return false
}


func Same [T comparable] (firstValue T, secondValue T) bool {
  return firstValue == secondValue
}
