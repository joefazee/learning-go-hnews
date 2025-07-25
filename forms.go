package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type formErrors map[string][]string

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (e formErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e formErrors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

type Form struct {
	url.Values
	Errors formErrors
}

func NewForm(form url.Values) *Form {
	return &Form{
		form,
		formErrors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) *Form {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("This field %s is required", field))
		}
	}
	return f
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MaxLength(field string, n int) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if utf8.RuneCountInString(value) > n {
		f.Errors.Add(field, fmt.Sprintf("This field  %s is too long (maximum of %d characters)", field, n))
	}

	return f
}

func (f *Form) MinLength(field string, n int) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if utf8.RuneCountInString(value) < n {
		f.Errors.Add(field, fmt.Sprintf("This field %s is too short (minimum of %d characters)", field, n))
	}

	return f
}

func (f *Form) Matches(field string, pattern *regexp.Regexp) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, fmt.Sprintf("This field %s is invalid", field))
	}

	return f
}

func (f *Form) IsEmail(field string) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}

	if !EmailRX.MatchString(value) {
		f.Errors.Add(field, fmt.Sprintf("This field %s is not a valid email address", field))
	}

	return f
}
