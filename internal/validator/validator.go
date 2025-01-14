package validator

import (
	"strings"
	"unicode/utf8"
)

// define a new validator type which contains a map of validation errors for our form fields
type Validator struct {
	FieldErrors map[string]string
}

// Valid() erturns true if the FieldErrors map doesnt contain ant entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError() adds and error message to the FieldErrors mao(so long as no entru already exists for the given key)
func (v *Validator) AddFieldError(key, message string) {
	// we need to initialize the map first, if it isnt already initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// ChechField() adds an error message to the FieldErrors map only if a validation check is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt() return true if a value is in a list of permitted intgers
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
