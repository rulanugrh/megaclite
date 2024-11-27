package middleware

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type IValidation interface {
	Validate(data any) error
	ValidationMessage(err error) error
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() IValidation {
	return &Validation{
		validate: validator.New(),
	}
}

func (v *Validation) Validate(data any) error {
	err := v.validate.Struct(data)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validation) ValidationMessage(err error) error {
	var msg string
	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			msg = fmt.Sprintf("This field is required: %v", e.Field())
		case "email":
			msg = fmt.Sprintf("This field must be email format: %s", e.Field())
		}
	}

	return web.InternalServerError(msg)
}
