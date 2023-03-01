package middlewares

import (
	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/gofiber/fiber/v2"
)

type AuthenticationMiddleware struct {
	validator *dependencies.Validator
}

func NewAuthenticationMiddleware(validator *dependencies.Validator) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		validator: validator,
	}
}

func (a *AuthenticationMiddleware) validateUserAndSetUserInHeader(c *fiber.Ctx) (entities.UserRead, error) {
	currentUser, err := helper.ValidateUserCredentical(c)
	if err != nil {
		return currentUser, err
	}

	c.Locals("currentUser", currentUser)

	return currentUser, nil
}

func (a *AuthenticationMiddleware) ValidateUser(c *fiber.Ctx) error {
	_, err := a.validateUserAndSetUserInHeader(c)
	if err != nil {
		return err
	}

	return c.Next()
}

func (a *AuthenticationMiddleware) ValidateAdmin(c *fiber.Ctx) error {
	currentUser, err := a.validateUserAndSetUserInHeader(c)
	if err != nil {
		return err
	}

	if !currentUser.IsAdmin {
		return c.Status(403).SendString("Only admin can do this action")
	}

	return c.Next()
}

func (a *AuthenticationMiddleware) ValidateUserSameAsUrlIdOrAdmin(c *fiber.Ctx) error {
	id, err := a.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	currentUser, err := a.validateUserAndSetUserInHeader(c)
	if err != nil {
		return err
	}

	if !currentUser.IsAdmin && currentUser.IdUser != id {
		return c.Status(403).SendString("Can't do this action to another user's account")
	}

	return c.Next()
}
