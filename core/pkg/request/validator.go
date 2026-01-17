package request

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/response"
	"backend/core/pkg/scope"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	Errors  map[string]string `json:"errors"`
	Code    int               `json:"code"`
	Success bool              `json:"success"`
}

type Validator struct {
	scope *scope.Scope
}

func NewValidator(scope *scope.Scope) *Validator {
	return &Validator{
		scope: scope,
	}
}

func (v *Validator) Validate(c *fiber.Ctx, model any) (bool, error) {
	if err := v.parse(c, model); err != nil {
		return false, response.NotValidRequest(c, v.scope, err)
	}

	if ok := v.validation(c, model); !ok {
		return false, nil
	}

	return true, nil
}

func (v *Validator) parse(c *fiber.Ctx, model any) error {
	var err error

	if c.Method() == http.MethodGet {
		err = c.QueryParser(model)
	} else {
		err = c.BodyParser(model)
	}

	if err == nil {
		return nil
	}

	c.Status(http.StatusBadRequest)
	return errorsx.Wrap(err, "Error parsing request body")
}

func (v *Validator) validation(c *fiber.Ctx, model any) bool {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validationErr := validate.Struct(model)
	if validationErr != nil {
		errs := make(map[string]string)
		for _, e := range validationErr.(validator.ValidationErrors) {
			errs[v.getJSONFieldName(e, model)] = v.validationMessage(e)
		}

		err := c.Status(http.StatusUnprocessableEntity).JSON(ValidationError{
			Errors:  errs,
			Code:    http.StatusUnprocessableEntity,
			Success: false,
		})
		if err != nil {
			return false
		}

		return false
	}

	return true
}

func (v *Validator) getJSONFieldName(fe validator.FieldError, model any) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fe.StructField()); ok {
		jsonTag := field.Tag.Get("json")
		return strings.Split(jsonTag, ",")[0]
	}

	return fe.Field()
}

func (v *Validator) validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "oneof":
		return fmt.Sprintf("Must be one of [%s]", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s'", fe.Tag())
	}
}
