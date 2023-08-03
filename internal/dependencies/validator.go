package dependencies

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	Validate *validator.Validate
}

func NewValidator(validate *validator.Validate) *Validator {
	return &Validator{Validate: validate}
}

func (v *Validator) formaFieldErrorMessage(fe validator.FieldError) string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + fe.Field() + "'")
	sb.WriteString(", condition: " + fe.ActualTag())

	// Print condition parameters, e.g. oneof=red blue -> { red blue }
	if fe.Param() != "" {
		sb.WriteString(" { " + fe.Param() + " }")
	}

	if fe.Value() != nil && fe.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", fe.Value()))
	}

	return sb.String()
}

func (v *Validator) validateStruct(payload interface{}) error {

	err := v.Validate.Struct(payload)
	errMessage := ""
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage += v.formaFieldErrorMessage(err) + "\n"
		}
		return fiber.NewError(400, errMessage)
	} else {
		return nil
	}
}

func (v *Validator) validateParse(c *fiber.Ctx, payload interface{}) error {
	err := v.validateStruct(payload)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return nil
}

func (v *Validator) ParseQuery(c *fiber.Ctx, queryStruct interface{}) error {
	err := c.QueryParser(queryStruct)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return v.validateParse(c, queryStruct)
}

func (v *Validator) ParseBody(c *fiber.Ctx, bodyStruct interface{}) error {
	v.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return name
	})

	err := c.BodyParser(bodyStruct)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return v.validateParse(c, bodyStruct)
}

func (v *Validator) ParseIdFromUrlParameter(c *fiber.Ctx) (int, error) {
	potentialId := c.Locals("id")
	if potentialId == nil {
		param := c.Params("id")
		err := v.Validate.Var(param, "required,number")
		if err != nil {
			return 0, fiber.NewError(400, "id parameter must be a valid positive integer")
		}

		id, err := strconv.Atoi(param)
		if err != nil {
			return 0, err
		}

		c.Locals("id", id)
		return id, err
	} else {
		id, ok := potentialId.(int)
		if !ok {
			return 0, errors.New("error, can't convert to int, variable error")
		}

		return id, nil
	}
}

func (v *Validator) GetAuthentication(c *fiber.Ctx) (entities.UserRead, error) {
	potentialUser := c.Locals("currentUser")
	if potentialUser == nil {
		panic(errors.New("error, there is no user set, please use the authentication middleware first"))
	}

	user, ok := potentialUser.(entities.UserRead)
	if !ok {
		panic(errors.New("error, can't convert to UserRead, variable error"))
	}

	return user, nil
}
