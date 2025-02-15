package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// use the regexp.MustCompile function to parse a regular expression pattern for sanity checking the format of an email address. this returns a pointer to a compiled regexp.Regexp type, or panics in the event of an error. parsing this pattern once at startup and stroing the compiled *regexp.Regexp in a variable is more preformant than reparsing the pattern each time we need
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)")

// define a new validator type which contains a map of validation errors for our form fields
// add a new NonFieldErrors []string field to the struct, which we will use to hold any validation errors which are not related to a specific form field
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() erturns true if the FieldErrors map doesnt contain ant entries
// update the valid() method to also check that the NonFieldErrors slice is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// create a AddNonFieldErrors helper for adding error mesaages to the new NonFieldErrors slice
func (v *Validator) AddNonFieldErrors(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
