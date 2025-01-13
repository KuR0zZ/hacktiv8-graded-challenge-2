package middleware

import (
	"graded-challenge-2-client/helper"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidate(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{validator}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat("invalid validation input"))
		}

		errMap := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errMap[field] = "this field is required"
			default:
				errMap[field] = "validation failed"
			}
		}

		return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat(errMap))
	}

	return nil
}
