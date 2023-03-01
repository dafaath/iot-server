package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dafaath/iot-server/v2/configs"
	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/dafaath/iot-server/v2/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	db         *pgxpool.Pool
	repository *repositories.UserRepository
	validator  *dependencies.Validator
}

func NewUserHandler(db *pgxpool.Pool, userRepository *repositories.UserRepository, validator *dependencies.Validator) (UserHandler, error) {
	return UserHandler{
		db:         db,
		validator:  validator,
		repository: userRepository,
	}, nil
}

func (u *UserHandler) Register(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	bodyPayload := entities.UserCreate{}

	err = u.validator.ParseBody(c, &bodyPayload)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	_, err = u.repository.GetByUsername(ctx, tx, bodyPayload.Username)
	if err != nil && !helper.IsErrorNotFound(err) {
		return err
	} else if err == nil {
		return fiber.NewError(fiber.StatusConflict, "Username already used")
	}

	_, err = u.repository.GetByEmail(ctx, tx, bodyPayload.Email)
	if err != nil && !helper.IsErrorNotFound(err) {
		return err
	} else if err == nil {
		return fiber.NewError(fiber.StatusConflict, "Email already used")
	}

	user, err := u.repository.Create(ctx, tx, bodyPayload)
	if err != nil {
		return err
	}

	// Untuk kepentingan testing, agar test otomatis tidak mengirim email
	sendEmail, err := strconv.ParseBool(c.Query("sendEmail", "true"))
	if err != nil {
		return err
	}

	if sendEmail {
		err = u.repository.SendEmailActivation(ctx, user)
		if err != nil {
			return err
		}
	}

	// Send token if env is test. For testing purpose.
	response := fmt.Sprintf("Success sign in, id: %d. Check email for activation", user.IdUser)
	config := configs.GetConfig()
	if config.Server.Env == "test" {
		jwt, err := u.repository.SignJWT(ctx, user)
		if err != nil {
			return err
		}
		response += fmt.Sprintf(". Token: %s|", jwt)
	}

	return c.Status(fiber.StatusCreated).SendString(response)
}

func (u *UserHandler) Login(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	bodyPayload := new(entities.UserLogin)
	err = u.validator.ParseBody(c, bodyPayload)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	user, err := u.repository.GetByUsername(ctx, tx, bodyPayload.Username)
	if err != nil {
		return helper.ChangeErrorIfErrorIsNotFound(err, fiber.NewError(401, "Username or password is incorrect"))
	}

	if !user.Status {
		return fiber.NewError(400, "Account is inactive, check email for activation")
	}

	err = u.repository.MatchPassword(ctx, tx, user, bodyPayload.Password)
	if err != nil {
		return fiber.NewError(401, "Username or password is incorrect")
	}

	token, err := u.repository.SignJWT(ctx, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(token)
}

func (u *UserHandler) Activation(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	query := new(entities.UserValidate)
	err = u.validator.ParseQuery(c, query)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	user, err := helper.ValidateUserToken(query.Token)
	if err != nil {
		return err
	}

	if user.Status {
		return fiber.NewError(400, "Your account has already activated")
	}

	err = u.repository.UpdateStatus(ctx, tx, user.IdUser, true)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Account for username: %s has been activated", user.Username))
}

func (u *UserHandler) ForgotPassword(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	body := new(entities.UserForgotPassword)
	err = u.validator.ParseBody(c, body)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	user, err := u.repository.GetByUsername(ctx, tx, body.Username)
	if err != nil {
		return helper.ChangeErrorIfErrorIsNotFound(err, fiber.NewError(400, "Username or email is incorrect"))
	}

	err = u.validator.Validate.VarWithValue(user.Email, body.Email, "eqcsfield")
	if err != nil {
		return fiber.NewError(400, "Username or email is incorrect")
	}

	if !user.Status {
		return fiber.NewError(403, "Your account is inactive. Check your email for activation")
	}

	newPassword := helper.GenerateRandomString(8)
	err = u.repository.UpdatePassword(ctx, tx, user.IdUser, newPassword)
	if err != nil {
		return err
	}

	err = u.repository.SendEmailForgotPassword(ctx, user, newPassword)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString("New password request sent. Check email for new password")
}

func (u *UserHandler) GetAll(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	users, err := u.repository.GetAll(ctx, tx)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (u *UserHandler) GetOne(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := u.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	user, err := u.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (u *UserHandler) Update(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := u.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	bodyPayload := new(entities.UserUpdatePassword)
	err = u.validator.ParseBody(c, bodyPayload)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	user, err := u.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = u.repository.MatchPassword(ctx, tx, user, bodyPayload.OldPassword)
	if err != nil {
		return fiber.NewError(401, "Old password is incorrect")
	}

	err = u.repository.UpdatePassword(ctx, tx, id, bodyPayload.NewPassword)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString("Success change password")
}

func (u *UserHandler) Delete(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := u.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	_, err = u.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = u.repository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Success delete user, id: %d", id))
}
