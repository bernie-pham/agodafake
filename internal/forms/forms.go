package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors{},
	}
}

func (f *Form) Has(field string, req *http.Request) bool {
	x := req.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// Valid return true if no error, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

func (f *Form) MinLength(field string, length int, req *http.Request) bool {
	x := req.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string, req *http.Request) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address")
		return false
	}
	return true
}

func (f *Form) IsEqual(field1, field2 string, req *http.Request) bool {
	value1 := req.Form.Get(field1)
	value2 := req.Form.Get(field2)
	if value1 != value2 {
		return false
	}
	return true
}
