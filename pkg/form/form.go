package form

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Form contains form data and any validation errs
type Form struct {
	ValidationErrs Errors
	Values         url.Values
	RequiredFields Required
}

// Required contains the max,min length for each required field
type Required map[string]struct {
	MinLength int
	MaxLength int
}

// Errors contains validation errs
type Errors map[string][]string

// Add adds and error to the Errors map
func (errs Errors) Add(field, err string) {
	errs[field] = append(errs[field], err)
}

// Get retrieves the errors for a given field
func (errs Errors) Get(field string) []string {
	return errs[field]
}

// New creates an instance of form
func (f *Form) New(data url.Values, required Required) *Form {
	return &Form{Values: data, RequiredFields: required, ValidationErrs: Errors{}}
}

// IsValid reports whether form contains any validation errs
func (f *Form) IsValid() bool {
	f.validateLength()
	f.validateRequiredFields()

	return len(f.ValidationErrs) == 0
}

// validateLength checks that a field is no longer than the permitted max length
func (f *Form) validateLength() {
	for fieldName, value := range f.RequiredFields {
		if utf8.RuneCountInString(f.Values.Get(fieldName)) > value.MaxLength {
			f.ValidationErrs.Add(fieldName, fmt.Sprintf("Too long! Keep it at %d characters max", value.MaxLength))
			return
		}

		if utf8.RuneCountInString(f.Values.Get(fieldName)) < value.MinLength {
			f.ValidationErrs.Add(fieldName, fmt.Sprintf("Too short! Use at least %d characters", value.MaxLength))
		}
	}
}

// validateRequiredFields checks that all required fields are present
// also checks that only [a-zA-Z0-9_] are present
func (f *Form) validateRequiredFields() {
	permissibleChars := regexp.MustCompile(`^[a-zA-Z0-9\ ]+$`)
	for field := range f.RequiredFields {
		value := f.Values.Get(field)

		if len(strings.Trim(value, "\t \n")) == 0 {
			f.ValidationErrs.Add(field, "Required!")
		}

		if !permissibleChars.MatchString(value) {
			f.ValidationErrs.Add(field, fmt.Sprintf("Invalid %s. Use characters(a-z) and numbers only", field))
		}
	}
}
