package validator

import (
	"regexp"
	"slices"
)

// Validator is a struct that holds validation errors
type Validator struct {
	Errors map[string]string
}

// New returns a new Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if there are no errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error to the validator
func (v *Validator) AddError(key string, err string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = err
	}
}

// Check checks if a condition is true
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// CheckFunc checks if a predicate function returns true
func (v *Validator) CheckFunc(predicate func() bool, key, message string) {
	v.Check(predicate(), key, message)
}

// Permitted checks if a value is in a list of permitted values
func Permitted[T comparable](value T, permitted ...T) bool {
	return slices.Contains(permitted, value)
}

// Matches checks if a value matches a pattern
func Matches(value string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(value)
}

// Unique checks if a slice of values is unique
func Unique[T comparable](values []T) bool {
	uniq := make(map[T]bool)

	for _, value := range values {
		uniq[value] = true
	}

	return len(uniq) == len(values)
}
