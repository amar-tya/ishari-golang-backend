package validation

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s any) error
}

type defaultValidator struct {
	v *validator.Validate
}

func New() Validator {
	return &defaultValidator{v: validator.New()}
}

func (d *defaultValidator) Struct(s any) error {
	return d.v.Struct(s)
}
