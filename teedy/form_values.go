package teedy

import (
	"fmt"
	"net/url"
)

type FormValues struct {
	values url.Values
	errors MultiError
}

func NewFormValues() *FormValues {
	return &FormValues{
		values: url.Values{},
		errors: MultiError{},
	}
}

func (f *FormValues) AddMandatory(key, value string) {
	if len(value) == 0 {
		f.errors[key] = fmt.Errorf("value for mandatory field '%s' is empty", key)
		return
	}

	f.values.Add(key, value)
}

func (f *FormValues) AddIfNotEmpty(key, value string) {
	if len(value) == 0 {
		return
	}

	f.values.Add(key, value)
}

func (f *FormValues) Result() (url.Values, error) {
	if len(f.errors) > 0 {
		return nil, f.errors
	}

	return f.values, nil
}

// MultiError stores multiple errors.
type MultiError map[string]error

func (e MultiError) Error() string {
	s := ""
	for _, err := range e {
		s = err.Error()
		break
	}
	switch len(e) {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, len(e)-1)
}

func (e MultiError) merge(errors MultiError) {
	for key, err := range errors {
		if e[key] == nil {
			e[key] = err
		}
	}
}
