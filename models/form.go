package models

import (
	"net/http"
	"strings"
)

// Form contains basic information about forms on html page.
type Form struct {
	InputNames []string
	Inputs     map[string]string
	Errors     map[string]string
}

// Populate extracts form values from http.Request
// and fills it's Inputs.
func (f *Form) Populate(r *http.Request) {
	for _, input := range f.InputNames {
		f.Inputs[input] = r.FormValue(input)
	}
}

// CheckRequiredFields ensures that required fields
// filled up properly
func (f *Form) CheckRequiredFields(fieldNames []string) {
	for _, field := range fieldNames {
		if strings.TrimSpace(f.Inputs[field]) == "" {
			f.Errors[field] = "Please enter a valid " + field + "."
		}
	}
}

// IsValid checks whether there is an error or not.
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// NewForm returns new form.
func NewForm(inputNames []string) *Form {
	return &Form{
		InputNames: inputNames,
		Inputs:     make(map[string]string),
		Errors:     make(map[string]string),
	}
}
