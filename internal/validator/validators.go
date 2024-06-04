package Validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) ==0
}

func (v *Validator) AddField(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exits := v.FieldErrors[key]; !exits {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddField(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func PermittedInt(value int, permitttedValues ...int) bool {
	for i := range permitttedValues {
		if value == i {
			return true
		}
	}
	return false
}
